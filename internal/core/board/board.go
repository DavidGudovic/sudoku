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
	CandidatePrefixRune = '*'
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
//   - Invalid is a board state where at least one row, column, or box contains duplicate values, or an EmptyCell contains NoCandidates
//   - Unsolved is a valid board state where at least one Cell is an EmptyCell
//   - Solved is a valid board state where all Cells are filled with values.
type State int

// Board represents a Size*Size Sudoku board composed of Cells.
// Each Cell holds its current value and a CandidateSet.
// The Board constraints (Sudoku rules) are always enforced, and all CandidateSet's are updated accordingly.
type Board struct {
	Cells [Size][Size]Cell
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

// filterCandidates removes candidate prefixes from the string representation of a board,
// returning only the values as a string.
// It returns an error if the string representation is invalid.
func filterCandidates(s string) (string, error) {
	var sb strings.Builder
	isCandidate := false

	for _, ch := range s {
		if isCandidate {
			isCandidate = false
			continue
		}

		if ch == CandidatePrefixRune {
			isCandidate = true
			continue
		}

		if ch < EmptyCell+'0' || ch > MaxValue+'0' {
			return "", ErrInvalidRuneInStringRep
		}

		sb.WriteRune(ch)
	}

	if isCandidate {
		return "", ErrInvalidStringRep
	}

	return sb.String(), nil
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
				for value := MinValue; value <= MaxValue; value++ {
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
// unless the value or the coordinates are illegal, in which case it returns a non-nil error.
func (b *Board) SetValueOnCoords(c Coordinates, value int) error {
	if value < EmptyCell || value > MaxValue {
		return ErrInvalidCellValue
	}

	if c.Row < 0 || c.Row >= Size || c.Col < 0 || c.Col >= Size {
		return ErrIndexOutOfBounds
	}

	b.Cells[c.Row][c.Col].value = value

	if value == EmptyCell {
		b.Cells[c.Row][c.Col].candidates = AllCandidates
		b.recalculateCandidateSets()
	}

	if value != EmptyCell {
		b.Cells[c.Row][c.Col].candidates = NoCandidates
		b.propagateConstraints(c, value)
	}

	return nil
}

// SetValueOnIndex sets the value on Coordinates corresponding to the given 0-based index,
// unless the value or the coordinates are illegal, in which case it returns a non-nil error.
func (b *Board) SetValueOnIndex(index int, value int) error {
	c, err := CoordsFromIndex(index)

	if err != nil {
		return err
	}

	return b.SetValueOnCoords(c, value)
}

// GetValueByIndex gets the value from Coordinates corresponding to the given 0-based index (left to right, top to bottom),
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
			c := NewCoordinates(row, col)
			cell := b.Cells[c.Row][c.Col]

			if cell.value == EmptyCell {
				if cell.candidates == NoCandidates {
					return Invalid
				}

				continue
			}

			if rows[c.Row].Contains(cell.value) || cols[c.Col].Contains(cell.value) || boxes[c.BoxIndex()].Contains(cell.value) {
				return Invalid
			}

			_ = rows[c.Row].Add(cell.value)
			_ = cols[c.Col].Add(cell.value)
			_ = boxes[c.BoxIndex()].Add(cell.value)
		}
	}

	return b.resolveValidState(rows, cols, boxes)
}

// resolveValidState is a helper function to check if every value has been seen in every row/col/box
// If it has, the Board is Solved, else it's Unsolved
func (b *Board) resolveValidState(rows, cols, boxes [Size]CandidateSet) State {
	for i := 0; i < Size; i++ {
		if (rows[i] & cols[i] & boxes[i]) != AllCandidates {
			return Unsolved
		}
	}

	return Solved
}

// propagateConstraints removes the given value from the CandidateSet of all Cells in the same row, column, and box as the given Coordinates
func (b *Board) propagateConstraints(c Coordinates, value int) {
	boxIndex := c.BoxIndex()

	for i := 0; i < Size; i++ {
		b.Cells[c.Row][i].candidates.Remove(value)
		b.Cells[i][c.Col].candidates.Remove(value)

		bc, _ := CoordsFromBoxIndex(boxIndex, i)
		b.Cells[bc.Row][bc.Col].candidates.Remove(value)
	}
}

// recalculateCandidateSets recalculates the CandidateSet for the Cells
func (b *Board) recalculateCandidateSets() {
	// Reset candidates for all empty cells
	for row := 0; row < Size; row++ {
		for col := 0; col < Size; col++ {
			if b.Cells[row][col].value == EmptyCell {
				b.Cells[row][col].candidates = AllCandidates
			}
		}
	}

	// Re-apply constraints based on current cell values
	for row := 0; row < Size; row++ {
		for col := 0; col < Size; col++ {
			cell := b.Cells[row][col]

			if cell.value != EmptyCell {
				b.propagateConstraints(NewCoordinates(row, col), cell.value)
			}
		}
	}
}
