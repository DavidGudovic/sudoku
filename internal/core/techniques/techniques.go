package techniques

import (
	"errors"

	"github.com/DavidGudovic/sudoku/internal/core/board"
)

var (
	ErrCannotProgress = errors.New("cannot progress in this step")
	ErrCannotSolve    = errors.New("cannot solve in this step")
)

// Technique represents a Sudoku solving technique that can be applied to a puzzle
type Technique interface {
	Apply(puzzle *board.Board) (Step, error)
}

// Func is a function type that implements the Technique interface
// This allows simple functions to be used as techniques
type Func func(*board.Board) (Step, error)

func (tf Func) Apply(puzzle *board.Board) (Step, error) {
	return tf(puzzle)
}

// Step represents a single solving step taken by a technique
type Step struct {
	Description       string
	Technique         string
	AffectedCells     []board.Coordinates
	ReasonCells       []board.Coordinates
	RemovedCandidates board.CandidateSet
	PlacedValue       *int
}

// StepStack represents a stack of Steps, typically used to track the sequence of solving steps
type StepStack []Step

func (s Step) MadeProgress() bool {
	return s.PlacedValue != nil || s.RemovedCandidates != board.NoCandidates
}

// rowCandidates returns a CandidateSet representing all candidates present in the given row at least once
func rowCandidates(p board.Board, row int) board.CandidateSet {
	seen := board.NoCandidates

	for col := 0; col < board.Size; col++ {
		cell := p.CellAt(board.Coordinates{Row: row, Col: col})
		seen.Merge(cell.Candidates())
	}

	return seen
}

// columnCandidates returns a CandidateSet representing all candidates present in the given column at least once
func columnCandidates(p board.Board, col int) board.CandidateSet {
	seen := board.NoCandidates

	for row := 0; row < board.Size; row++ {
		cell := p.CellAt(board.Coordinates{Row: row, Col: col})
		seen.Merge(cell.Candidates())
	}

	return seen
}

// boxCandidates returns a CandidateSet representing all candidates present in the given box at least once
func boxCandidates(p board.Board, coords board.Coordinates) board.CandidateSet {
	seen := board.NoCandidates

	boxIndex := coords.BoxIndex()

	for i := 0; i < board.BoxSize*board.BoxSize; i++ {
		bc, _ := board.CoordsFromBoxIndex(boxIndex, i)
		cell := p.CellAt(bc)
		seen.Merge(cell.Candidates())
	}

	return seen
}
