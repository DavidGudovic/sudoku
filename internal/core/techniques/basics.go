package techniques

import (
	"github.com/DavidGudovic/sudoku/internal/core/board"
)

// LastDigit technique:
//
// If a candidate can only fit in one cell of a row, column, or box, place it there.
// TODO: Refactor
func LastDigit(puzzle *board.Board) (Step, error) {
	var candidates board.CandidateSet
	var coords board.Coordinates
	var val int

	for row := 0; row < board.Size; row++ {
		for col := 0; col < board.Size; col++ {
			coords = board.MustCoordinates(row, col)
			candidates = puzzle.CellAt(coords).Candidates()

			if candidates.Count() == 1 {
				val = candidates.First()
				targetCellCandidates := board.MustCandidateSet(val)

				step := Step{
					Technique:         "LastDigit",
					AffectedCells:     NoPeers.With(coords),
					RemovedCandidates: targetCellCandidates,
					PlacedValue:       &val,
				}

				rowPeers := Peers.Of(coords).Across(Row)
				columnPeers := Peers.Of(coords).Across(Column)
				boxPeers := Peers.Of(coords).Across(Box)

				found := false

				if rowPeers.With(coords).Candidates(*puzzle) == targetCellCandidates {
					found = true
					step.Technique += " (Row)"
					step.ReasonCells = rowPeers
				} else if columnPeers.With(coords).Candidates(*puzzle) == targetCellCandidates {
					found = true
					step.Technique += " (Column)"
					step.ReasonCells = columnPeers
				} else if boxPeers.With(coords).Candidates(*puzzle) == targetCellCandidates {
					found = true
					step.Technique += " (Box)"
					step.ReasonCells = boxPeers
				}

				if found {
					return step.MustApplyTo(puzzle), nil
				}
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
