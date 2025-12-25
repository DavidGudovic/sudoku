package techniques

import (
	"errors"

	"github.com/DavidGudovic/sudoku/internal/core/board"
)

var (
	ErrCannotProgress     = errors.New("cannot progress in this step")
	ErrCannotSolve        = errors.New("cannot solve in this step")
	ErrUnapplicablePuzzle = errors.New("step cannot be applied to puzzle, wrong puzzle? ")
)

// Technique is a function type that implements a Sudoku solving technique.
type Technique func(*board.Board) (Step, error)

// Step represents a single solving step taken by a Technique
// It contains information about the technique used, affected cells, candidates removed, and values placed, useful for explanatory UI's
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

// MadeProgress returns true if the step resulted in any progress (either placing a value or removing candidates)
func (s Step) MadeProgress() bool {
	return s.PlacedValue != nil || s.RemovedCandidates != board.NoCandidates
}

// ApplyTo applies the step to the puzzle, returning an error if the step cannot be applied.
func (s Step) ApplyTo(puzzle *board.Board) error {
	for _, coords := range s.AffectedCells {
		if s.PlacedValue != nil {
			err := puzzle.SetValueOnCoords(coords, *s.PlacedValue)
			if err != nil {
				return ErrUnapplicablePuzzle
			}
		}

		if s.RemovedCandidates != board.NoCandidates {
			puzzle.ExcludeCandidatesFromCoords(coords, s.RemovedCandidates)
		}
	}

	return nil
}

// MustApplyTo applies the step to the puzzle, panicking if the step cannot be applied.
//
// Only to be used when the passed puzzle is the puzzle on which the step was generated.
// Panicking here indicates a serious error in technique code or world view, invalidating any further program flow.
func (s Step) MustApplyTo(puzzle *board.Board) Step {
	if err := s.ApplyTo(puzzle); err != nil {
		panic("Impossible: " + err.Error())
	}

	return s
}
