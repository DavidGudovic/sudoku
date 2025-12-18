package techniques

import (
	"fmt"

	"github.com/DavidGudovic/sudoku/internal/core/board"
)

// HiddenSingle technique:
// For any value If in any row or box only once cell has the value as a candidate, the value must be there.
// This is true even when the cell has more than one candidate itself.
func HiddenSingle(puzzle *board.Board) (Step, error) {
	p := Peers.All()

	for coords := range p.Populated() {
		candidates := puzzle.CellAt(coords).Candidates()

		if candidates.Count() < 2 {
			continue
		}

		var scope Scope
		var scopeIndex int

		for _, candidate := range candidates.Slice() {
			boxPeers := Peers.Of(coords).Across(Box).ContainingCandidates(*puzzle, board.MustCandidateSet(candidate))
			rowPeers := Peers.Of(coords).Across(Row).ContainingCandidates(*puzzle, board.MustCandidateSet(candidate))
			colPeers := Peers.Of(coords).Across(Column).ContainingCandidates(*puzzle, board.MustCandidateSet(candidate))

			if boxPeers == NoPeers {
				scope = Box
				scopeIndex = coords.BoxIndex()
			} else if rowPeers == NoPeers {
				scope = Row
				scopeIndex = coords.Row
			} else if colPeers == NoPeers {
				scope = Column
				scopeIndex = coords.Col
			} else {
				continue
			}

			step := Step{
				Technique:         "HiddenSingle (" + scope.String() + ")",
				AffectedCells:     []board.Coordinates{coords},
				ReasonCells:       Peers.Of(coords).Across(scope).EmptyCells(*puzzle).Slice(),
				Description:       fmt.Sprint("None of the empty cells in ", scope, " ", scopeIndex, " can hold a ", candidate, " except ", coords, ", placing a ", candidate),
				RemovedCandidates: board.MustCandidateSet(candidate),
				PlacedValue:       &candidate,
			}

			return step.MustApplyTo(puzzle), nil
		}
	}

	return Step{}, ErrCannotProgress
}

func HiddenPair(_ *board.Board) (Step, error) {
	return Step{}, nil
}
