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
