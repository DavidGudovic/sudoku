package board

import (
	"errors"
	"strings"
)

const (
	Size            = 9
	BoxSize         = 3
	BoxCount        = 9
	CellCount       = 81
	Invalid   State = "Invalid"
	Unsolved  State = "Unsolved"
	Solved    State = "Solved"
)

var (
	ErrInvalidStringRep       = errors.New("invalid string representation")
	ErrInvalidRuneInStringRep = errors.New("invalid rune in string")
	ErrIndexOutOfBounds       = errors.New("index out of bounds")
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
				for _, value := range AllCellValues {
					if cell.candidates.Contains(value) {
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

func (b *Board) GetState() State {
	var rows [Size]BitMask
	var cols [Size]BitMask
	var boxes [BoxCount]BitMask

	for row := 0; row < Size; row++ {
		for col := 0; col < Size; col++ {
			value := b.Cells[row][col].value

			if value == EmptyCell {
				continue
			}

			box := (row/BoxSize)*BoxSize + (col / BoxSize)

			if rows[row].Contains(value) || cols[col].Contains(value) || boxes[box].Contains(value) {
				return Invalid
			}

			rows[row].Add(value)
			cols[col].Add(value)
			boxes[box].Add(value)
		}
	}

	return b.resolveState(rows, cols, boxes)
}

func (b *Board) resolveState(rows, cols, boxes [Size]BitMask) State {
	for i := 0; i < Size; i++ {
		if (rows[i] & cols[i] & boxes[i]) != AllCandidates {
			return Unsolved
		}
	}

	return Solved
}
