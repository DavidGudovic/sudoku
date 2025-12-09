package board

import (
	"errors"
	"strings"
)

const (
	Size      = 9
	BoxSize   = 3
	CellCount = Size * Size
	Invalid   = "Invalid"
	Unsolved  = "Unsolved"
	Solved    = "Solved"
)

var (
	ErrInvalidStringRep       = errors.New("invalid string representation")
	ErrInvalidRuneInStringRep = errors.New("invalid rune in string")
	ErrIndexOutOfBounds       = errors.New("index out of bounds")
	ErrInvalidBoardState      = errors.New("invalid board state")
)

type State string
type Board struct {
	Cells [Size][Size]Cell
}

func NewBoard() *Board {
	cells := [Size][Size]Cell{}

	for row := 0; row < Size; row++ {
		for col := 0; col < Size; col++ {
			cells[row][col] = NewCell()
		}
	}

	return &Board{Cells: cells}
}

func FromString(s string) (*Board, error) {
	board := NewBoard()
	valuesOnly, err := filterCandidates(s)

	if err != nil {
		return nil, err
	}

	if len(valuesOnly) != CellCount {
		return nil, ErrInvalidStringRep
	}

	for i := 0; i < CellCount; i++ {
		row, col, err := CoordsFromIndex(i)

		if err != nil {
			return nil, err
		}

		err = board.SetValueOnCoords(row, col, int(valuesOnly[i]-'0'))

		if err != nil {
			return nil, err
		}
	}

	return board, nil
}

func (b *Board) ToString(withCandidates bool) string {
	var s strings.Builder

	for row := 0; row < Size; row++ {
		for col := 0; col < Size; col++ {
			cell := b.Cells[row][col]

			s.WriteRune(rune(cell.value + '0'))

			if withCandidates && cell.value == EmptyCell {
				for value := range AllCellValues {
					if cell.candidates.Includes(value) {
						s.WriteRune(CandidatePrefixRune)
						s.WriteRune(rune(value + '0'))
					}
				}
			}
		}
	}

	return s.String()
}

func (b *Board) SetValueOnCoords(row, col, value int) error {
	if value < EmptyCell || value > MaxValue {
		return ErrInvalidCellValue
	}

	if row < 0 || row >= Size || col < 0 || col >= Size {
		return ErrIndexOutOfBounds
	}

	b.Cells[row][col].value = value

	if value != EmptyCell {
		b.Cells[row][col].candidates = NoCandidates
	} else {
		b.Cells[row][col].candidates = AllCandidates
	}

	return nil
}

func (b *Board) SetValueOnIndex(index int, value int) error {
	row, col, err := CoordsFromIndex(index)

	if err != nil {
		return err
	}

	return b.SetValueOnCoords(row, col, value)
}

func (b *Board) ValidateState() State {
	allRowsSolved, err := validateRows(b)
	if err != nil {
		return Invalid
	}

	allColsSolved, err := validateCols(b)
	if err != nil {
		return Invalid
	}

	allBoxesSolved, err := validateBoxes(b)
	if err != nil {
		return Invalid
	}

	solved := allRowsSolved && allColsSolved && allBoxesSolved

	if solved {
		return Solved
	}

	return Unsolved
}

func validateRows(b *Board) (bool, error) {
	var seen [MaxValue + 1]bool
	allRowsSolved := true

	for _, rowArray := range b.Cells {
		for _, cell := range rowArray {
			if cell.value == EmptyCell {
				continue
			}

			if seen[cell.value] {
				return false, ErrInvalidBoardState
			}

			seen[cell.value] = true
		}
		allRowsSolved = allRowsSolved && allValuesSeen(seen)
		seen = [MaxValue + 1]bool{}
	}

	return allRowsSolved, nil
}

func validateCols(b *Board) (bool, error) {
	var seen [MaxValue + 1]bool
	allColsSolved := true

	for col := 0; col < Size; col++ {
		for row := 0; row < Size; row++ {
			value := b.Cells[row][col].value

			if value == EmptyCell {
				continue
			}

			if seen[value] {
				return false, ErrInvalidBoardState
			}

			seen[value] = true
		}

		allColsSolved = allColsSolved && allValuesSeen(seen)
		seen = [MaxValue + 1]bool{}
	}

	return allColsSolved, nil
}

func validateBoxes(b *Board) (bool, error) {
	var seen [MaxValue + 1]bool
	allBoxesSolved := true

	for boxRow := range BoxSize {
		for boxCol := range BoxSize {

			for rowInBox := 0; rowInBox < BoxSize; rowInBox++ {
				for colInBox := 0; colInBox < BoxSize; colInBox++ {
					rowInSudoku := boxRow*BoxSize + rowInBox
					colInSudoku := boxCol*BoxSize + colInBox

					value := b.Cells[rowInSudoku][colInSudoku].value

					if value == EmptyCell {
						continue
					}

					if seen[value] {
						return false, ErrInvalidBoardState
					}

					seen[value] = true
				}
			}

			allBoxesSolved = allBoxesSolved && allValuesSeen(seen)
			seen = [MaxValue + 1]bool{}
		}
	}

	return allBoxesSolved, nil
}

func allValuesSeen(seen [MaxValue + 1]bool) bool {
	for _, value := range AllCellValues {
		if !seen[value] {
			return false
		}
	}

	return true
}
