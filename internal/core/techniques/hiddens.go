package techniques

import (
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
						AffectedCells:     NoPeers.With(coords),
						ReasonCells:       scopedPeers.EmptyCells(*puzzle),
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

// HiddenPair technique:
//
// If two candidates can only fit in the same two cells of a Scope, and those cells contain other candidates as well,
// those other candidates can be removed from those cells.
func HiddenPair(puzzle *board.Board) (Step, error) {
	return hiddenMultiple(puzzle, 2, "HiddenPair")
}

// HiddenTriple technique:
//
// If three candidates can only fit in the same three cells of a Scope, and those cells contain other candidates as well,
// those other candidates can be removed from those cells.
func HiddenTriple(puzzle *board.Board) (Step, error) {
	return hiddenMultiple(puzzle, 3, "HiddenTriple")
}

// HiddenQuad technique:
//
// If four candidates can only fit in the same four cells of a Scope, and those cells contain other candidates as well,
// those other candidates can be removed from those cells.
func HiddenQuad(puzzle *board.Board) (Step, error) {
	return hiddenMultiple(puzzle, 4, "HiddenQuad")
}

// hiddenMultiple is a helper function for HiddenPair, HiddenTriple, and HiddenQuad techniques.
func hiddenMultiple(puzzle *board.Board, count int, techniqueName string) (Step, error) {
	for i := 0; i < board.Size; i++ {
		for _, scope := range AllScopes {
			inScope := Peers.InScope(scope, i).EmptyCells(*puzzle)

			if inScope.Count() < count {
				continue
			}

			peerCandidates := inScope.Candidates(*puzzle)

			if peerCandidates.Count() <= count { // No hidden multiples possible, if its equal its naked, not hidden
				continue
			}

			combinations := board.CandidateSubsets(peerCandidates, count)

			for _, combo := range combinations {
				potentialMultiples := inScope.ContainingCandidates(*puzzle, combo)

				if potentialMultiples.Count() != count {
					continue
				}

				affected := potentialMultiples.Candidates(*puzzle)

				affected.Exclude(combo)

				if affected == board.NoCandidates {
					continue
				}

				step := Step{
					Technique:         techniqueName,
					AffectedCells:     potentialMultiples,
					ReasonCells:       potentialMultiples,
					RemovedCandidates: affected,
				}

				return step.MustApplyTo(puzzle), nil
			}
		}
	}

	return Step{}, ErrCannotProgress
}
