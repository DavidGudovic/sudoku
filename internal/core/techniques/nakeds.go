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
			coords := board.MustCoordinates(row, col)
			candidates := puzzle.CellAt(coords).Candidates()

			if candidates.Count() == 1 {
				val := candidates.First()
				step := Step{
					Technique:         "NakedSingle",
					AffectedCells:     []board.Coordinates{coords},
					ReasonCells:       Peers.Of(coords).Across(AllScopes...).Slice(),
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
// If two cells in a Scope contain the same pair of candidates,
// those candidates can be removed from all other cells in that Scope containing them.
func NakedPair(puzzle *board.Board) (Step, error) {
	for row := 0; row < board.Size; row++ {
		for col := 0; col < board.Size; col++ {
			candidates := puzzle.CellAt(board.MustCoordinates(row, col)).Candidates()

			if candidates.Count() != 2 {
				continue
			}

			found := board.MustCoordinates(row, col)
			potentialPairs := Peers.Of(found).Across(AllScopes...).ContainingExactCandidates(*puzzle, candidates)

			if potentialPairs.IsEmpty() {
				continue
			}

			var affected PeerSet
			var pair board.Coordinates

			for peer := range potentialPairs.Each() {
				affected = Peers.Of(found, peer).AcrossSharedScopes().ContainingCandidates(*puzzle, candidates)

				if affected.IsEmpty() {
					continue
				}

				pair = peer
				break
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

// NakedTriple technique:
//
// If three cells in a Scope contain the same triplet of candidates,
// those candidates can be removed from all other cells in that Scope containing them.
func NakedTriple(puzzle *board.Board) (Step, error) {
	return Step{}, ErrCannotProgress
}
