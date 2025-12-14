package techniques

import (
	"fmt"

	"github.com/DavidGudovic/sudoku/internal/core/board"
)

// NakedSingle technique:
//
// If a cell has only one candidate left, place it there.
func NakedSingle(puzzle *board.Board) (Step, error) {
	for row := 0; row < board.Size; row++ {
		for col := 0; col < board.Size; col++ {
			coords, _ := board.NewCoordinates(row, col)
			cell := puzzle.CellAt(coords)
			candidates := cell.Candidates()

			if candidates.Count() == 1 {
				val := candidates.Slice()[0]
				step := Step{
					Technique:         "NakedSingle",
					AffectedCells:     []board.Coordinates{coords},
					ReasonCells:       AllPeersOf(coords).Slice(),
					RemovedCandidates: candidates,
					PlacedValue:       &val,
					Description:       fmt.Sprint("The candidate ", val, " is the only one left at ", coords, ", placing a ", val),
				}

				err := puzzle.SetValueOnCoords(coords, val)

				if err != nil {
					panic("Impossible: " + err.Error())
				}

				return step, nil
			}
		}
	}

	return Step{}, ErrCannotProgress
}

// NakedPair technique:
//
// If two cells in a group (row, column, or box) contain the same pair of candidates,
// those candidates can be removed from all other cells in that group containing them.
func NakedPair(puzzle *board.Board) (Step, error) {
	for row := 0; row < board.Size; row++ {
		for col := 0; col < board.Size; col++ {
			foundCoords, _ := board.NewCoordinates(row, col)

			candidates := puzzle.CellAt(foundCoords).Candidates()

			if candidates.Count() != 2 {
				continue
			}

			peers := AllPeersOf(foundCoords).ContainingExactCandidates(*puzzle, candidates)

			if peers.IsEmpty() {
				continue
			}

			pairCoords := peers.Slice()[0]

			if InSameRow(foundCoords, pairCoords) {
				rowPeers := RowPeersOf(foundCoords).NotContainingCandidates(*puzzle, candidates)
			} else if InSameColumn(foundCoords, pairCoords) {
				rowPeers := ColumnPeersOf(foundCoords).NotContainingCandidates(*puzzle, candidates)
			} else if InSameBox(foundCoords, pairCoords) {
				rowPeers := BoxPeersOf(foundCoords).NotContainingCandidates(*puzzle, candidates)
			}
		}
	}

	return Step{}, nil
}
