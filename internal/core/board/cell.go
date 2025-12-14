package board

import (
	"errors"
	"fmt"
)

const (
	EmptyCell = 0
	MinValue  = 1
	MaxValue  = 9
)

var ErrInvalidCellValue = errors.New("invalid cell value")

// Coordinates represents the row and column of a cell on the Board.
type Coordinates struct {
	Row int
	Col int
}

// Cell represents a single cell on the Sudoku board.
// It holds the current value EmptyCell or MinValue -> MaxValue and a CandidateSet
type Cell struct {
	value      int
	candidates CandidateSet
}

// NewCoordinates creates a new Coordinates struct
func NewCoordinates(row, col int) (Coordinates, error) {
	if row < 0 || row >= Size || col < 0 || col >= Size {
		return Coordinates{}, ErrIndexOutOfBounds
	}

	return Coordinates{
		Row: row,
		Col: col,
	}, nil
}

// CoordsFromIndex converts a 0-based index (left to right, top to bottom) to Coordinates.
func CoordsFromIndex(index int) (Coordinates, error) {
	if index < 0 || index >= CellCount {
		return Coordinates{}, ErrIndexOutOfBounds
	}

	c, _ := NewCoordinates(index/Size, index%Size)
	return c, nil
}

// CoordsFromBoxIndex converts a box index (0-8) and position within the box (0-8) to Coordinates.
func CoordsFromBoxIndex(boxIndex int, positionInBox int) (Coordinates, error) {
	if boxIndex < 0 || boxIndex >= BoxCount {
		return Coordinates{}, ErrIndexOutOfBounds
	}

	if positionInBox < 0 || positionInBox >= BoxSize*BoxSize {
		return Coordinates{}, ErrIndexOutOfBounds
	}

	row := (boxIndex/BoxSize)*BoxSize + (positionInBox / BoxSize)
	col := (boxIndex%BoxSize)*BoxSize + (positionInBox % BoxSize)

	c, _ := NewCoordinates(row, col)

	return c, nil
}

// NewCell creates a new Cell with an EmptyCell value and AllCandidates available.
func NewCell() Cell {
	return Cell{
		value:      EmptyCell,
		candidates: AllCandidates,
	}
}

// BoxIndex returns the box index (0-8) of the Coordinates.
func (c Coordinates) BoxIndex() int {
	return (c.Row/BoxSize)*BoxSize + c.Col/BoxSize
}

func (c Coordinates) String() string {
	return fmt.Sprintf("R%dC%d", c.Row, c.Col)
}

func (c Cell) Value() int {
	return c.value
}

func (c Cell) Candidates() CandidateSet {
	return c.candidates
}

func (c Cell) ContainsCandidate(val int) bool {
	return c.candidates.Contains(val)
}
