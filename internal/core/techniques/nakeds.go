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
			candidates := puzzle.CellAt(board.Coordinates{Row: row, Col: col}).Candidates()

			if candidates.Count() != 2 {
				continue
			}

			found, _ := board.NewCoordinates(row, col)
			peers := AllPeersOf(found).ContainingExactCandidates(*puzzle, candidates)

			if peers.IsEmpty() {
				continue
			}

			var affected PeerSet
			var pair board.Coordinates

			for _, peer := range peers.Slice() {
				if found.SharesRowWith(peer) {
					affected = RowPeersOf(found).ContainingCandidates(*puzzle, candidates).Excluding(found, peer)
					pair = peer
				} else if found.SharesColumnWith(peer) {
					affected = ColumnPeersOf(found).ContainingCandidates(*puzzle, candidates).Excluding(found, peer)
					pair = peer
				} else if found.SharesBoxWith(peer) {
					affected = BoxPeersOf(found).ContainingCandidates(*puzzle, candidates).Excluding(found, peer)
					pair = peer
				}
			}

			if affected.IsEmpty() {
				continue
			}

			step := Step{
				Technique:         "NakedPair",
				Description:       fmt.Sprint("Naked Pair found at ", found, " and ", pair, ", removing candidates from peers"),
				AffectedCells:     []board.Coordinates{found, pair},
				ReasonCells:       affected.Slice(),
				RemovedCandidates: candidates,
			}

			for _, coords := range affected.Slice() {
				puzzle.ExcludeCandidatesFromCoords(coords, candidates)
			}

			return step, nil
		}
	}

	return Step{}, ErrCannotProgress
}
