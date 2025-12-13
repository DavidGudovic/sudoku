package techniques

import "github.com/DavidGudovic/sudoku/internal/core/board"

func (pp PointingPair) Apply(_ *board.Board) (Step, error) {
	return Step{}, nil
}

func (ld LastDigit) Apply(_ *board.Board) (Step, error) {
	return Step{}, nil
}
