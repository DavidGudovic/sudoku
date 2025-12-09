package solver

import (
	"github.com/DavidGudovic/sudoku/internal/core/board"
)

type Solver interface {
	Solve(b *board.Board) error
}

type Validator interface {
	ValidateBoardState(b *board.Board) (bool, error)
}

type BacktrackingSolver struct {
	Validator Validator
}

type StrategySolver struct {
	Validator Validator
}

func NewBacktrackingSolver(v Validator) *BacktrackingSolver {
	return &BacktrackingSolver{Validator: v}
}

func NewStrategySolver(v Validator) *StrategySolver {
	return &StrategySolver{Validator: v}
}

func (s *BacktrackingSolver) Solve(b *board.Board) error {
	return nil
}

func (s *StrategySolver) Solve(b *board.Board) error {
	return nil
}
