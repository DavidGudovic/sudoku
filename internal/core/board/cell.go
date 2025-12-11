package board

import "errors"

const (
	EmptyCell = 0
	MinValue  = 1
	MaxValue  = 9
)

var (
	// AllCandidates 0b1111111110 BitMask.
	AllCandidates = BitMask(1022)
	// NoCandidates 0b0000000000 BitMask.
	NoCandidates = BitMask(0)

	AllCellValues = [9]int{
		1, 2, 3, 4, 5, 6, 7, 8, 9,
	}

	ErrInvalidCellValue = errors.New("invalid cell value")
)

// BitMask represents a set of candidates for a cell using bitwise operations.
// The least significant bit (bit 0) is unused, bits 1-9 represent candidates 1-9.
// For example, a BitMask of 0b0000100110 represents candidates 1, 2, and 5.
type BitMask uint16

type Cell struct {
	value      int
	candidates BitMask
}

// NewCell creates a new Cell with an EmptyCell value and AllCandidates available.
func NewCell() Cell {
	return Cell{
		value:      EmptyCell,
		candidates: AllCandidates,
	}
}

// Contains checks if the BitMask contains the specified candidate value.
func (bm *BitMask) Contains(value int) bool {
	return *bm&(1<<value) != 0
}

// Add adds (sets the bit to 1) the specified candidate value to the BitMask.
func (bm *BitMask) Add(value int) {
	*bm |= 1 << value
}

// Remove removes (sets the bit to 0) the specified candidate value from the BitMask.
func (bm *BitMask) Remove(value int) {
	*bm &= ^(1 << value)
}

func (c Cell) Value() int {
	return c.value
}

func (c Cell) Candidates() BitMask {
	return c.candidates
}
