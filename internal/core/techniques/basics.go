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
				val = candidates.ToSlice()[0]
				targetCellCandidates, _ := board.NewCandidateSet(val)

				step := Step{
					Technique:         "LastDigit",
					AffectedCells:     []board.Coordinates{coords},
					RemovedCandidates: targetCellCandidates,
					PlacedValue:       &val,
				}

				rowPeers := RowPeersOf(coords)
				columnPeers := ColumnPeersOf(coords)
				boxPeers := BoxPeersOf(coords)

				if rowPeers.With(coords).Candidates(*puzzle) == targetCellCandidates {
					step.Technique += " (Row)"
					step.Description = fmt.Sprint("Value ", val, " can only go in one place in Row ", coords.Row, ", placing a ", val, " at ", coords)
					step.ReasonCells = rowPeers.Slice()
					return step, nil
				}
				if columnPeers.With(coords).Candidates(*puzzle) == targetCellCandidates {
					step.Technique += " (Column)"
					step.Description = fmt.Sprint("Value ", val, " can only go in one place in Col ", coords.Col, ", placing a ", val, " at ", coords)
					step.ReasonCells = columnPeers.Slice()
					return step, nil
				}
				if boxPeers.With(coords).Candidates(*puzzle) == targetCellCandidates {
					step.Technique += " (Box)"
					step.Description = fmt.Sprint("Value ", val, " can only go in one place in Box ", coords.BoxIndex(), ", placing a ", val, " at ", coords)
					step.ReasonCells = boxPeers.Slice()
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
