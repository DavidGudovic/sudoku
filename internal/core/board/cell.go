package board

import "errors"

const (
	CandidatePrefixRune = '_'
	EmptyCell           = 0
	MinValue            = 1
	MaxValue            = 9
)

var (
	AllCandidates = BitMask(1022)
	NoCandidates  = BitMask(0)

	AllCellValues = [9]int{
		1, 2, 3, 4, 5, 6, 7, 8, 9,
	}

	ErrInvalidCellValue = errors.New("invalid cell value")
)

type BitMask uint16

type Cell struct {
	value      int
	candidates BitMask
}

func NewCell() Cell {
	return Cell{
		value:      EmptyCell,
		candidates: AllCandidates,
	}
}

func (bm *BitMask) Contains(value int) bool {
	return *bm&(1<<value) != 0
}

func (bm *BitMask) Add(value int) {
	*bm |= 1 << value
}

func (bm *BitMask) Remove(value int) {
	*bm &= ^(1 << value)
}

func (c Cell) Value() int {
	return c.value
}

func (c Cell) Candidates() BitMask {
	return c.candidates
}
