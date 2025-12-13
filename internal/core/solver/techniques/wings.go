package techniques

import "github.com/DavidGudovic/sudoku/internal/core/board"

func (xw XWing) Apply(_ *board.Board) (Step, error) {
	return Step{}, nil
}

func (ss Skyscraper) Apply(_ *board.Board) (Step, error) {
	return Step{}, nil
}

func (tsk TwoStringKite) Apply(_ *board.Board) (Step, error) {
	return Step{}, nil
}
