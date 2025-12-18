package techniques

import (
	"fmt"

	"github.com/DavidGudovic/sudoku/internal/core/board"
)

// HiddenSingle technique:
//
// For any value If in any row or box only one cell has the value as a candidate, the value must be there.
// This is true even when the cell has more than one candidate itself.
func HiddenSingle(puzzle *board.Board) (Step, error) {
	scopes := [3]struct {
		scope Scope
		index int
	}{
		{Row, 0},
		{Column, 0},
		{Box, 0},
	}

	for coords := range Peers.All().Each() {
		candidates := puzzle.CellAt(coords).Candidates()

		if candidates.Count() < 2 { // Cannot be hidden without at least 2 candidates
			continue
		}

		scopes[0].index = coords.Row
		scopes[1].index = coords.Col
		scopes[2].index = coords.BoxIndex()

		for candidate := range candidates.Each() {
			candidateMask := board.MustCandidateSet(candidate)

			for _, s := range scopes {
				scopedPeers := Peers.Of(coords).Across(s.scope)

				if scopedPeers.ContainingCandidates(*puzzle, candidateMask) == NoPeers {
					step := Step{
						Technique:         "HiddenSingle (" + s.scope.String() + ")",
						AffectedCells:     []board.Coordinates{coords},
						ReasonCells:       scopedPeers.EmptyCells(*puzzle).Slice(),
						Description:       fmt.Sprint("None of the empty cells in ", s.scope, " ", s.index, " can hold a ", candidate, " except ", coords, ", placing a ", candidate),
						RemovedCandidates: candidateMask,
						PlacedValue:       &candidate,
					}

					return step.MustApplyTo(puzzle), nil
				}
			}

		}
	}

	return Step{}, ErrCannotProgress
}

func HiddenPair(_ *board.Board) (Step, error) {
	return Step{}, nil
}
