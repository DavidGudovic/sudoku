package techniques

import (
	"fmt"

	"github.com/DavidGudovic/sudoku/internal/core/board"
)

// LastDigit technique:
//
// If a candidate can only fit in one cell of a row, column, or box, place it there.
func LastDigit(puzzle *board.Board) (Step, error) {
	var candidates board.CandidateSet
	var coords board.Coordinates
	var val int

	for row := 0; row < board.Size; row++ {
		for col := 0; col < board.Size; col++ {
			coords, _ = board.NewCoordinates(row, col)
			candidates = puzzle.CellAt(coords).Candidates()

			if candidates.Count() == 1 {
				rCandidates := rowCandidates(*puzzle, row)
				cCandidates := columnCandidates(*puzzle, col)
				bCandidates := boxCandidates(*puzzle, coords)

				val = candidates.ToSlice()[0]
				targetCellCandidates, _ := board.NewCandidateSet(val)

				step := Step{
					Technique:         "LastDigit",
					AffectedCells:     []board.Coordinates{coords},
					RemovedCandidates: targetCellCandidates,
					PlacedValue:       &val,
				}

				if rCandidates == targetCellCandidates {
					step.Technique = step.Technique + " (Row)"
					step.Description = fmt.Sprint("Value ", val, " can only go in one place in Row ", coords.Row, ", placing a ", val, " at ", coords)
					step.ReasonCells = rowPeers(coords).Slice()
					return step, nil
				}
				if cCandidates == targetCellCandidates {
					step.Technique = step.Technique + " (Column)"
					step.Description = fmt.Sprint("Value ", val, " can only go in one place in Col ", coords.Col, ", placing a ", val, " at ", coords)
					step.ReasonCells = colPeers(coords).Slice()
					return step, nil
				}
				if bCandidates == targetCellCandidates {
					step.Technique = step.Technique + " (Box)"
					step.Description = fmt.Sprint("Value ", val, " can only go in one place in Box ", coords.BoxIndex(), ", placing a ", val, " at ", coords)
					step.ReasonCells = boxPeers(coords).Slice()
					return step, nil
				}
			}
		}
	}

	return Step{}, ErrCannotProgress
}

func PointingPair(_ *board.Board) (Step, error) {
	return Step{}, nil
}
