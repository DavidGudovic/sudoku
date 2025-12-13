package techniques

import (
	"github.com/DavidGudovic/sudoku/internal/core/board"
)

type Technique interface {
	Apply(puzzle *board.Board) (Step, error)
}

type Step struct {
	Description       string
	Technique         string
	AffectedCells     []board.Coordinates
	ReasonCells       []board.Coordinates
	RemovedCandidates board.CandidateSet
	PlacedValue       *int
}

type StepStack []Step

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

func (ns NakedSingle) Apply(_ *board.Board) (Step, error) {
	return Step{}, nil
}

func (hs HiddenSingle) Apply(_ *board.Board) (Step, error) {
	return Step{}, nil
}

func (np NakedPair) Apply(_ *board.Board) (Step, error) {
	return Step{}, nil
}

func (hp HiddenPair) Apply(_ *board.Board) (Step, error) {
	return Step{}, nil
}

func (pp PointingPair) Apply(_ *board.Board) (Step, error) {
	return Step{}, nil
}

func (xw XWing) Apply(_ *board.Board) (Step, error) {
	return Step{}, nil
}

func (ss Skyscraper) Apply(_ *board.Board) (Step, error) {
	return Step{}, nil
}
