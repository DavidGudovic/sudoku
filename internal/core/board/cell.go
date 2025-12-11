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

	AllCellValues = [9]int{
		1, 2, 3, 4, 5, 6, 7, 8, 9,
	}

	ErrInvalidCellValue = errors.New("invalid cell value")
)

// CandidateSet is a bit mask representing the possible candidates for a cell.
// The least significant bit (bit 0) is unused, bits 1-9 represent candidates 1-9.
//   - For example, a CandidateSet of 0b0000100110 represents candidates 1, 2, and 5.
type CandidateSet uint16

type Cell struct {
	value      int
	candidates CandidateSet
}

// NewCell creates a new Cell with an EmptyCell value and AllCandidates available.
func NewCell() Cell {
	return Cell{
		value:      EmptyCell,
		candidates: AllCandidates,
	}
}

// Contains checks if the CandidateSet contains the specified candidate value.
func (bm *CandidateSet) Contains(value int) bool {
	return *bm&(1<<value) != 0
}

// Add adds (sets the bit to 1) the specified candidate value to the CandidateSet.
func (bm *CandidateSet) Add(value int) error {
	if value < MinValue || value > MaxValue {
		return ErrInvalidCellValue
	}

	*bm |= 1 << value

	return nil
}

// Remove removes (sets the bit to 0) the specified candidate value from the CandidateSet.
func (bm *CandidateSet) Remove(value int) {
	*bm &= ^(1 << value)
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
