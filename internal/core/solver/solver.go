package solver

import (
	"github.com/DavidGudovic/sudoku/internal/core/board"
)

type Solver interface {
	Solve(b *board.Board) error
}

type BacktrackingSolver struct{}

type StrategySolver struct{}

func (s *BacktrackingSolver) Solve(_ *board.Board) error {
	return nil
}

func (s *StrategySolver) Solve(_ *board.Board) error {
	return nil
}
