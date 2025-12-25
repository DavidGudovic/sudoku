package techniques

import (
	"github.com/DavidGudovic/sudoku/internal/core/board"
)

// LastDigit technique:
//
// If a candidate can only fit in one cell of a row, column, or box, place it there.
func LastDigit(puzzle *board.Board) (Step, error) {
	relevant := Peers.All().EmptyCells(*puzzle).ContainingCountCandidates(*puzzle, 1)

	if relevant.IsEmpty() {
		return Step{}, ErrCannotProgress
	}

	for c := range relevant.Each() {
		candidates := puzzle.CellAt(c).Candidates()
		val := candidates.First()

		for _, s := range AllScopes {
			scopedPeers := Peers.Of(c).Across(s)

			if scopedPeers.ContainingCandidates(*puzzle, candidates).IsEmpty() {
				step := Step{
					Technique:         "LastDigit (" + s.String() + ")",
					AffectedCells:     NoPeers.With(c),
					ReasonCells:       scopedPeers,
					RemovedCandidates: candidates,
					PlacedValue:       &val,
				}

				return step.MustApplyTo(puzzle), nil
			}
		}
	}

	return Step{}, ErrCannotProgress
}

// LockedCandidates technique:
//
// If candidates are confined to a single scope within a box,
// those candidates can be removed from the shared scopes of the cells containing them.
func LockedCandidates(puzzle *board.Board) (Step, error) {
	relevant := Peers.All().EmptyCells(*puzzle)
	candidates := relevant.Candidates(*puzzle)

	for candidate := range candidates.Each() {
		for _, s := range [2]Scope{Row, Column} {
			for i := 0; i < board.Size; i++ {
				mask := board.MustCandidateSet(candidate)
				peers := Peers.InScope(s, i).ContainingCandidates(*puzzle, mask)

				if peers.Count() > 3 || peers.Count() < 2 {
					continue
				}

				if peers.ShareScope(Box) == false {
					continue
				}

				affected := Peers.Of(peers.Slice()...).AcrossSharedScopes().ContainingCandidates(*puzzle, mask)

				if affected.IsEmpty() {
					continue
				}

				step := Step{
					Technique:         "LockedCandidates (" + s.String() + ")",
					AffectedCells:     affected,
					ReasonCells:       peers,
					RemovedCandidates: mask,
				}

				return step.MustApplyTo(puzzle), nil
			}
		}
	}

	return Step{}, ErrCannotProgress
}
