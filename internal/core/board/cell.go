package board

import "errors"

const (
	EmptyCell = 0
	MinValue  = 1
	MaxValue  = 9
)

var (
	AllCandidates CandidateSet = 0b1111111110
	NoCandidates  CandidateSet = 0b0000000000

	ErrInvalidCellValue = errors.New("invalid cell value")
)

// Coordinates represents the row and column of a cell on the Board.
type Coordinates struct {
	Row int
	Col int
}

// CandidateSet is a bit mask representing the possible candidates for a Cell.
// The least significant bit (bit 0) is unused, bits 1-9 represent candidates 1-9.
//   - For example, a CandidateSet of 0b0000100110 represents candidates 1, 2, and 5.
type CandidateSet uint16

// Cell represents a single cell on the Sudoku board.
// It holds the current value EmptyCell or MinValue -> MaxValue and a CandidateSet
type Cell struct {
	value      int
	candidates CandidateSet
}

// NewCoordinates creates a new Coordinates struct
func NewCoordinates(row, col int) Coordinates {
	return Coordinates{
		Row: row,
		Col: col,
	}
}

// CoordsFromIndex converts a 0-based index (left to right, top to bottom) to Coordinates.
func CoordsFromIndex(index int) (Coordinates, error) {
	if index < 0 || index >= CellCount {
		return NewCoordinates(0, 0), ErrIndexOutOfBounds
	}

	return NewCoordinates(index/Size, index%Size), nil
}

// CoordsFromBoxIndex converts a box index (0-8) and position within the box (0-8) to Coordinates.
func CoordsFromBoxIndex(boxIndex int, positionInBox int) (Coordinates, error) {
	if boxIndex < 0 || boxIndex >= BoxCount {
		return NewCoordinates(0, 0), ErrIndexOutOfBounds
	}

	if positionInBox < 0 || positionInBox >= BoxSize*BoxSize {
		return NewCoordinates(0, 0), ErrIndexOutOfBounds
	}

	row := (boxIndex/BoxSize)*BoxSize + (positionInBox / BoxSize)
	col := (boxIndex%BoxSize)*BoxSize + (positionInBox % BoxSize)

	return NewCoordinates(row, col), nil
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
	return "(" + string(rune(c.Row+'0')) + ", " + string(rune(c.Col+'0')) + ")"
}

// Contains checks if the CandidateSet contains the specified candidate value.
func (cs *CandidateSet) Contains(value int) bool {
	return *cs&(1<<value) != 0
}

// Add adds (sets the bit to 1) the specified candidate value to the CandidateSet.
func (cs *CandidateSet) Add(value int) error {
	if value < MinValue || value > MaxValue {
		return ErrInvalidCellValue
	}

	*cs |= 1 << value

	return nil
}

// Remove removes (sets the bit to 0) the specified candidate value from the CandidateSet.
func (cs *CandidateSet) Remove(value int) {
	*cs &= ^(1 << value)
}

// Exclude removes all candidates present in the other CandidateSet from the current CandidateSet.
func (cs *CandidateSet) Exclude(other CandidateSet) {
	*cs &^= other
}

// Merge adds all candidates present in the other CandidateSet to the current CandidateSet.
func (cs *CandidateSet) Merge(other CandidateSet) {
	*cs |= other
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
