package techniques

import "github.com/DavidGudovic/sudoku/internal/core/board"

func (ns NakedSingle) Apply(_ *board.Board) (Step, error) {
	return Step{}, nil
}

func (np NakedPair) Apply(_ *board.Board) (Step, error) {
	return Step{}, nil
}
