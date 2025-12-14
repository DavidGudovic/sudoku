package techniques

import (
	"fmt"

	"github.com/DavidGudovic/sudoku/internal/core/board"
)

func NakedSingle(puzzle *board.Board) (Step, error) {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			coords, _ := board.NewCoordinates(row, col)
			cell := puzzle.CellAt(coords)
			candidates := cell.Candidates()

			if candidates.Count() == 1 {
				val := candidates.ToSlice()[0]
				step := Step{
					Technique:         "NakedSingle",
					AffectedCells:     []board.Coordinates{coords},
					ReasonCells:       AllPeersOf(coords).Slice(),
					RemovedCandidates: candidates,
					PlacedValue:       &val,
					Description:       fmt.Sprint("The candidate ", val, " is the only one left at ", coords, ", placing a ", val),
				}
				return step, nil
			}
		}
	}

	return Step{}, ErrCannotProgress
}

func NakedPair(_ *board.Board) (Step, error) {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			_ = board.Coordinates{}
			_ = board.NoCandidates
		}
	}
	return Step{}, nil
}
