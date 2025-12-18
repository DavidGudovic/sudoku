package techniques

import "github.com/DavidGudovic/sudoku/internal/core/board"

// HiddenSingle technique:
// For any value If in any row or box only once cell has the value as a candidate, the value must be there.
// This is true even when the cell has more than one candidate itself.
func HiddenSingle(puzzle *board.Board) (Step, error) {
	p := Peers.All()

	step := Step{
		Technique: "HiddenSingle",
	}

	// For every coordinate
	for coords := range p.Populated() {
		candidates := puzzle.CellAt(coords).Candidates()

		// If it has less than 2 candidates, it can't be hidden'
		if candidates.Count() < 2 {
			continue
		}

		var scope Scope

		// For every candidate of the current cell
		for _, candidate := range candidates.Slice() {
			boxPeers := Peers.Of(coords).Across(Box).ContainingCandidates(*puzzle, board.MustCandidateSet(candidate))
			rowPeers := Peers.Of(coords).Across(Row).ContainingCandidates(*puzzle, board.MustCandidateSet(candidate))
			colPeers := Peers.Of(coords).Across(Column).ContainingCandidates(*puzzle, board.MustCandidateSet(candidate))

			if boxPeers == NoPeers {
				scope = Box
			} else if rowPeers == NoPeers {
				scope = Row
			} else if colPeers == NoPeers {
				scope = Column
			} else {
				continue
			}

			step = Step{
				AffectedCells:     []board.Coordinates{coords},
				ReasonCells:       Peers.Of(coords).Across(scope).Slice(),
				Description:       "In " + scope.String() + ", value " + candidate.String() + " can only go in one place at " + coords.String() + ", placing a " + candidate.String(),
				RemovedCandidates: board.MustCandidateSet(candidate),
				PlacedValue:       &candidate,
			}

			step.Technique += " " + scope.String()

			return step.MustApplyTo(puzzle), nil
		}
	}

	return step, ErrCannotProgress
}

func HiddenPair(_ *board.Board) (Step, error) {
	return Step{}, nil
}
