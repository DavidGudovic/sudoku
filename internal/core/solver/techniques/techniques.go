package techniques

import (
	"github.com/DavidGudovic/sudoku/internal/core/board"
)

type Technique interface {
	Apply(puzzle board.Board) (Step, board.Board, error)
	ExpectsToSolve() bool
}

type Step struct {
	Description       string
	Technique         string
	AffectedCells     []board.Coordinates
	ReasonCells       []board.Coordinates
	RemovedCandidates board.CandidateSet
	PlacedValue       *int
}

func (s Step) MadeProgress() bool {
	return s.PlacedValue != nil || s.RemovedCandidates != board.NoCandidates
}

type (
	NakedSingle  struct{}
	HiddenSingle struct{}
	NakedPair    struct{}
	HiddenPair   struct{}
	PointingPair struct{}
	XWing        struct{}
	Skyscraper   struct{}
)

func (ns NakedSingle) Apply(puzzle board.Board) (Step, board.Board, error) {
	return Step{}, puzzle, nil
}

func (ns NakedSingle) ExpectsToSolve() bool {
	return false
}

func (hs HiddenSingle) Apply(puzzle board.Board) (Step, board.Board, error) {
	return Step{}, puzzle, nil
}

func (hs HiddenSingle) ExpectsToSolve() bool {
	return false
}

func (np NakedPair) Apply(puzzle board.Board) (Step, board.Board, error) {
	return Step{}, puzzle, nil
}

func (np NakedPair) ExpectsToSolve() bool {
	return false
}

func (hp HiddenPair) Apply(puzzle board.Board) (Step, board.Board, error) {
	return Step{}, puzzle, nil
}

func (hp HiddenPair) ExpectsToSolve() bool {
	return false
}

func (pp PointingPair) Apply(puzzle board.Board) (Step, board.Board, error) {
	return Step{}, puzzle, nil
}

func (pp PointingPair) ExpectsToSolve() bool {
	return false
}

func (xw XWing) Apply(puzzle board.Board) (Step, board.Board, error) {
	return Step{}, puzzle, nil
}

func (xw XWing) ExpectsToSolve() bool {
	return false
}

func (ss Skyscraper) Apply(puzzle board.Board) (Step, board.Board, error) {
	return Step{}, puzzle, nil
}

func (ss Skyscraper) ExpectsToSolve() bool {
	return false
}
