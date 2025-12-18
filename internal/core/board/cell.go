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

// Coordinates represents the row and column of a Cell on the Board.
type Coordinates struct {
	Row int
	Col int
}

// Cell represents a single Cell on the Sudoku board.
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

// MustCoordinates creates a new Coordinates struct assuming the caller guarantees valid input
// If the input is invalid, there's either a serious bug in the caller or the world view is wrong, therefore it panics
func MustCoordinates(row, col int) Coordinates {
	c, err := NewCoordinates(row, col)
	if err != nil {
		panic("Impossible: " + err.Error())
	}

	return c
}

// CoordsFromIndex converts a 0-based index (left to right, top to bottom) to Coordinates.
func CoordsFromIndex(index int) (Coordinates, error) {
	if index < 0 || index >= CellCount {
		return Coordinates{}, ErrIndexOutOfBounds
	}

	return MustCoordinates(index/Size, index%Size), nil
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

	return MustCoordinates(row, col), nil
}

// MustCoordsFromBoxIndex is like CoordsFromBoxIndex but assumes the caller guarantees valid input.
// If the input is invalid, there's either a serious bug in the caller or the world view is wrong, therefore it panics.
func MustCoordsFromBoxIndex(boxIndex int, positionInBox int) Coordinates {
	c, err := CoordsFromBoxIndex(boxIndex, positionInBox)

	if err != nil {
		panic("Impossible: " + err.Error())
	}

	return c
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

func (c Coordinates) SharesRowWith(other Coordinates) bool {
	return c.Row == other.Row
}

func (c Coordinates) SharesColumnWith(other Coordinates) bool {
	return c.Col == other.Col
}

func (c Coordinates) SharesBoxWith(other Coordinates) bool {
	return c.BoxIndex() == other.BoxIndex()
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

func (c Cell) IsEmpty() bool {
	return c.value == EmptyCell
}

func (c Cell) ContainsCandidate(val int) bool {
	return c.candidates.Contains(val)
}
