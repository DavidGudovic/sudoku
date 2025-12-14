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

				return step.MustApplyTo(puzzle), nil
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
				}

				if found.SharesBoxWith(peer) {
					affected = affected.Union(BoxPeersOf(found).ContainingCandidates(*puzzle, candidates).Excluding(found, peer))
					pair = peer
				}

				if !affected.IsEmpty() {
					break
				}
			}

			if affected.IsEmpty() {
				continue
			}

			step := Step{
				Technique:         "NakedPair",
				Description:       fmt.Sprint("Naked Pair found at ", found, " and ", pair, ", removing candidates ", candidates.String(), " from peers"),
				AffectedCells:     affected.Slice(),
				ReasonCells:       []board.Coordinates{found, pair},
				RemovedCandidates: candidates,
			}

			return step.MustApplyTo(puzzle), nil
		}
	}

	return Step{}, ErrCannotProgress
}
