package techniques

import (
	"errors"

	"github.com/DavidGudovic/sudoku/internal/core/board"
)

type Backtracking struct{}

// Apply attempts to solve the puzzle using backtracking.
// It does not produce incremental steps, but returns the solved puzzle if successful.
func (bt Backtracking) Apply(puzzle board.Board) (Step, board.Board, error) {
	solvedPuzzle, err := bt.backtrackSolve(puzzle)

	if err != nil {
		return Step{}, puzzle, err
	}

	return Step{
		Description:       "Backtracking does not produce incremental steps, but can solve any puzzle.",
		Technique:         "Backtracking",
		RemovedCandidates: board.AllCandidates,
	}, solvedPuzzle, nil
}

func (bt Backtracking) ExpectsToSolve() bool {
	return true
}

// backtrackSolve implements the backtracking algorithm to solve the Sudoku puzzle.
// It recursively fills empty cells with candidate values and backtracks upon encountering invalid states.
// TODO optimize with constraint propagation
func (bt Backtracking) backtrackSolve(puzzle board.Board) (board.Board, error) {
	c, err := bt.findEmptyCell(puzzle)

	if err != nil {
		return puzzle, nil
	}

	for _, val := range board.AllCellValues {
		if puzzle.Cells[c.Row][c.Col].ContainsCandidate(val) == false {
			continue
		}
		_ = puzzle.SetValueOnCoords(c, val)

		switch puzzle.GetState() {
		case board.Invalid:
			_ = puzzle.SetValueOnCoords(c, board.EmptyCell)
			continue
		case board.Solved:
			return puzzle, nil
		default:
			solvedPuzzle, err := bt.backtrackSolve(puzzle)

			if err == nil {
				return solvedPuzzle, nil
			}

			_ = puzzle.SetValueOnCoords(c, board.EmptyCell)
		}
	}

	return puzzle, errors.New("no solution found")
}

// findEmptyCell is a helper that returns the coordinates of the first empty cell found in the puzzle.
// Searches left to right, top to bottom.
func (bt Backtracking) findEmptyCell(puzzle board.Board) (board.Coordinates, error) {
	for row := 0; row < board.Size; row++ {
		for col := 0; col < board.Size; col++ {
			if puzzle.Cells[row][col].Value() == board.EmptyCell {
				return board.Coordinates{Row: row, Col: col}, nil
			}
		}
	}

	return board.Coordinates{}, errors.New("no empty cells found")
}
