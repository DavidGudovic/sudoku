package board

//go:generate stringer -type=State

import (
	"errors"
	"strings"
)

const (
	Size                = 9
	BoxSize             = 3
	BoxCount            = 9
	CellCount           = 81
	CandidatePrefixRune = '_'
)

const (
	Invalid State = iota
	Unsolved
	Solved
)

var (
	ErrInvalidStringRep       = errors.New("invalid string representation")
	ErrInvalidRuneInStringRep = errors.New("invalid rune in string")
	ErrIndexOutOfBounds       = errors.New("index out of bounds")
)

// State represents the current state of the board
//   - Invalid is a board state in which any row/column or a box has > 1 of any value.
//   - Unsolved is a board state which != Invalid, but still has Cells with no value (EmptyCell value)
//   - Solved is a board state in which all Cells have values, and no row/column or box has > 1 of any value.
type State int

type Coordinates struct {
	Row int
	Col int
}

type Board struct {
	Cells [Size][Size]Cell
}

// NewCoordinates creates a new Coordinates struct
func NewCoordinates(row, col int) Coordinates {
	return Coordinates{Row: row, Col: col}
}

// NewBoard initializes a Size*Size empty board with full candidates
func NewBoard() *Board {
	cells := [Size][Size]Cell{}

	for row := 0; row < Size; row++ {
		for col := 0; col < Size; col++ {
			cells[row][col] = NewCell()
		}
	}

	return &Board{Cells: cells}
}

// FromString populates a new empty board with the values extracted from a string representation of the board.
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
		c, err := CoordsFromIndex(i)

		if err != nil {
			return nil, err
		}

		err = board.SetValueOnCoords(c, int(valuesOnly[i]-'0'))

		if err != nil {
			return nil, err
		}
	}

	return board, nil
}

// ToString extracts a string representation from the current state of a board.
// The Board is read left to right, top to bottom, where each value is represented as is, and the candidates are prefixed with CandidatePrefixRune.
func (b *Board) ToString(withCandidates bool) string {
	var s strings.Builder

	for row := 0; row < Size; row++ {
		for col := 0; col < Size; col++ {
			cell := b.Cells[row][col]

			s.WriteRune(rune(cell.value + '0'))

			if withCandidates && cell.value == EmptyCell {
				for _, value := range AllCellValues {
					if cell.ContainsCandidate(value) {
						s.WriteRune(CandidatePrefixRune)
						s.WriteRune(rune(value + '0'))
					}
				}
			}
		}
	}

	return s.String()
}

// SetValueOnCoords sets the value on Coordinates provided,
// unless the value, or the coordinates are illegal, in which case it returns a non-nil error.
func (b *Board) SetValueOnCoords(c Coordinates, value int) error {
	if value < EmptyCell || value > MaxValue {
		return ErrInvalidCellValue
	}

	if c.Row < 0 || c.Row >= Size || c.Col < 0 || c.Col >= Size {
		return ErrIndexOutOfBounds
	}

	b.Cells[c.Row][c.Col].value = value

	if value != EmptyCell {
		b.Cells[c.Row][c.Col].candidates = NoCandidates
	} else {
		b.Cells[c.Row][c.Col].candidates = AllCandidates
	}

	return nil
}

// SetValueOnIndex sets the value on Coordinates corresponding to the given 0-based index,
// unless the value, or the coordinates are illegal, in which case it returns a non-nil error.
func (b *Board) SetValueOnIndex(index int, value int) error {
	c, err := CoordsFromIndex(index)

	if err != nil {
		return err
	}

	return b.SetValueOnCoords(c, value)
}

// GetValueByIndex gets the value from Coordinates corresponding to the given 0-based index,
// unless the index is illegal, in which case it returns ErrIndexOutOfBounds
func (b *Board) GetValueByIndex(index int) (int, error) {
	c, err := CoordsFromIndex(index)

	if err != nil {
		return 0, err
	}

	return b.Cells[c.Row][c.Col].value, nil
}

// GetState resolves the current state of the board into one of [ Invalid, Unsolved, Solved]
// See [State] for more information.
func (b *Board) GetState() State {
	var rows [Size]CandidateSet
	var cols [Size]CandidateSet
	var boxes [BoxCount]CandidateSet

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

			_ = rows[row].Add(value)
			_ = cols[col].Add(value)
			_ = boxes[box].Add(value)
		}
	}

	return b.resolveState(rows, cols, boxes)
}

// resolveState is a helper function to check if every value has been seen in every row/col/box
// If it has, the Board is Solved, else it's Unsolved
func (b *Board) resolveState(rows, cols, boxes [Size]CandidateSet) State {
	for i := 0; i < Size; i++ {
		if (rows[i] & cols[i] & boxes[i]) != AllCandidates {
			return Unsolved
		}
	}

	return Solved
}
