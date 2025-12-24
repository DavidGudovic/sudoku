package techniques

import (
	"fmt"

	"github.com/DavidGudovic/sudoku/internal/core/board"
)

// NakedSingle technique:
//
// If a cell has only one candidate left, place it there.
func NakedSingle(puzzle *board.Board) (Step, error) {
	for c := range Peers.All().EmptyCells(*puzzle).ContainingCountCandidates(*puzzle, 1).Each() {
		candidates := puzzle.CellAt(c).Candidates()
		val := candidates.First()

		step := Step{
			Technique:         "NakedSingle",
			AffectedCells:     []board.Coordinates{c},
			ReasonCells:       Peers.Of(c).Across(AllScopes...).Slice(),
			RemovedCandidates: candidates,
			PlacedValue:       &val,
			Description:       fmt.Sprint("The candidate ", val, " is the only one left at ", c, ", placing a ", val),
		}

		return step.MustApplyTo(puzzle), nil
	}

	return Step{}, ErrCannotProgress
}

// NakedPair technique:
//
// If two cells in a Scope contain the same pair of candidates,
// those candidates can be removed from all other cells in that Scope containing them.
func NakedPair(puzzle *board.Board) (Step, error) {
	return nakedMultiple(puzzle, 2, "NakedPair")
}

// NakedTriple technique:
//
// If three cells in a Scope contain the same triplet of candidates, or a subset of them,
// those candidates can be removed from all other cells in that Scope containing them.
func NakedTriple(puzzle *board.Board) (Step, error) {
	return nakedMultiple(puzzle, 3, "NakedTriple")
}

// NakedQuad technique:
//
// If four cells in a Scope contain the same four candidates, or a subset of them,
// those candidates can be removed from all other cells in that Scope containing them.
func NakedQuad(puzzle *board.Board) (Step, error) {
	return nakedMultiple(puzzle, 4, "NakedQuad")
}

// nakedMultiple is a helper function for NakedPair, NakedTriple, and NakedQuad techniques.
func nakedMultiple(puzzle *board.Board, count int, techniqueName string) (Step, error) {
	for i := 0; i < board.Size; i++ {
		for _, scope := range AllScopes {
			potentialMultiples := Peers.InScope(scope, i).EmptyCells(*puzzle).ContainingMaxCandidates(*puzzle, count)

			if potentialMultiples.Count() < count {
				continue
			}

			combinations := Subsets(potentialMultiples, count)

			for _, combo := range combinations {
				candidates := combo.Candidates(*puzzle)

				if candidates.Count() == count {
					affected := Peers.Of(combo.Slice()...).AcrossSharedScopes().
						ContainingCandidates(*puzzle, candidates).
						Excluding(combo.Slice()...)

					if affected.IsEmpty() {
						continue
					}

					step := Step{
						Technique:         techniqueName,
						Description:       fmt.Sprint(techniqueName, " found at ", combo.String(), ", removing candidates ", candidates.String(), " from mutual peers"),
						AffectedCells:     affected.Slice(),
						ReasonCells:       combo.Slice(),
						RemovedCandidates: candidates,
					}

					return step.MustApplyTo(puzzle), nil
				}
			}

		}
	}

	return Step{}, ErrCannotProgress
}
