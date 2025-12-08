package board

import (
	"errors"
	"strings"
)

const (
	One = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
)

const (
	Size                = 9
	RowCount            = Size
	ColCount            = Size
	BoxSize             = 3
	BoxRowCount         = BoxSize
	BoxColCount         = BoxSize
	MinValue            = One
	MaxValue            = Nine
	EmptyCell           = 0
	CellCount           = RowCount * ColCount
	BoxCellCount        = BoxRowCount * BoxColCount
	CandidatePrefixRune = '_'
)

var (
	AllCandidates = [10]bool{false, true, true, true, true, true, true, true, true, true}
	NoCandidates  = [10]bool{false, false, false, false, false, false, false, false, false}

	AllCellValues = [9]int{
		One, Two, Three, Four, Five, Six, Seven, Eight, Nine,
	}

	ErrInvalidStringRepresentation = errors.New("invalid string representation")
	ErrInvalidRuneInString         = errors.New("invalid rune in string")
	ErrIndexOutOfBounds            = errors.New("index out of bounds")
)

type Board struct {
	Cells [RowCount][ColCount]*Cell
}

type Cell struct {
	Value      int
	Candidates [10]bool
}

func NewCell() *Cell {
	return &Cell{
		Value:      EmptyCell,
		Candidates: AllCandidates,
	}
}

func NewBoard() *Board {
	cells := [RowCount][ColCount]*Cell{}

	for row := 0; row < RowCount; row++ {
		for col := 0; col < ColCount; col++ {
			cells[row][col] = NewCell()
		}
	}

	return &Board{Cells: cells}
}

func (b *Board) SetValue(row, col, value int) {
	b.Cells[row][col].Value = value
	b.Cells[row][col].Candidates = NoCandidates
}

func (b *Board) ToString(withCandidates bool) string {
	var s strings.Builder

	for row := 0; row < RowCount; row++ {
		for col := 0; col < ColCount; col++ {
			cell := b.Cells[row][col]

			s.WriteRune(rune(cell.Value + '0'))

			if withCandidates && cell.Value != EmptyCell {
				for value := range AllCellValues {
					if cell.Candidates[value] {
						s.WriteRune(CandidatePrefixRune)
						s.WriteRune(rune(value + '0'))
					}
				}
			}
		}
	}

	return s.String()
}

func FromString(s string) (*Board, error) {
	board := NewBoard()
	cleanString, err := extractValuesStringRepresentation(s)

	if err != nil {
		return nil, err
	}

	cleanStringLength := len(cleanString)

	if cleanStringLength != CellCount {
		return nil, ErrInvalidStringRepresentation
	}

	for i := 0; i < cleanStringLength; i++ {
		row, col, err := CoordsFromIndex(i)

		if err != nil {
			return nil, err
		}

		board.SetValue(row, col, int(s[i]-'0'))
	}

	return board, nil
}

func CoordsFromIndex(index int) (int, int, error) {
	if index < 0 || index >= CellCount {
		return 0, 0, ErrIndexOutOfBounds
	}

	return index / 9, index % 9, nil
}

func extractValuesStringRepresentation(s string) (string, error) {
	var sb strings.Builder
	isCandidate := false

	for _, ch := range s {
		if isCandidate {
			continue
		}

		if ch == CandidatePrefixRune {
			isCandidate = true
			continue
		}

		if ch < EmptyCell+'0' || ch > rune(MaxValue+'0') {
			return "", ErrInvalidRuneInString
		}

		sb.WriteRune(ch)
	}

	if isCandidate {
		return "", ErrInvalidStringRepresentation
	}

	return sb.String(), nil
}
