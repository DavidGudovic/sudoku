package techniques

import "github.com/DavidGudovic/sudoku/internal/core/board"

func (hs HiddenSingle) Apply(_ *board.Board) (Step, error) {
	return Step{}, nil
}

func (hp HiddenPair) Apply(_ *board.Board) (Step, error) {
	return Step{}, nil
}
