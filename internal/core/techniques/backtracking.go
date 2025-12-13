package techniques

import (
	"errors"
	"fmt"
	"time"

	"github.com/DavidGudovic/sudoku/internal/core/board"
)

var (
	ErrCannotSolve = errors.New("cannot backtrack further; puzzle is unsolvable")
)

type Backtracking struct{}

type backtrackStats struct {
	guesses    int
	backtracks int
}

// Apply attempts to solve the puzzle using backtracking.
// It does not produce incremental steps but returns the solved puzzle if successful.
func (bt Backtracking) Apply(puzzle *board.Board) (Step, error) {
	stats := backtrackStats{}

	if puzzle.IsFaux() && puzzle.GetState() == board.Invalid {
		return Step{}, ErrCannotSolve
	}

	start := time.Now()
	solvedPuzzle, err := bt.backtrackSolve(*puzzle, &stats)
	elapsed := time.Since(start)

	if err != nil {
		return Step{}, err
	}

	madeChanges := *puzzle != solvedPuzzle

	*puzzle = solvedPuzzle

	description := fmt.Sprintf("Guessed %d times.\nBacktracked %d times.\nSolved the puzzle in %d milliseconds.",
		stats.guesses, stats.backtracks, elapsed.Milliseconds())

	step := Step{
		Description: description,
		Technique:   "Backtracking",
	}

	if madeChanges {
		step.RemovedCandidates = board.AllCandidates
	}

	return step, nil
}

// backtrackSolve implements the backtracking algorithm to solve the Sudoku puzzle.
// It recursively fills empty cells with candidate values and backtracks upon encountering invalid states.
func (bt Backtracking) backtrackSolve(puzzle board.Board, stats *backtrackStats) (board.Board, error) {
	c, err := bt.findEmptyCell(puzzle)

	if err != nil {
		return puzzle, nil
	}

	for val := board.MinValue; val <= board.MaxValue; val++ {
		if !puzzle.IsFaux() && puzzle.Cells[c.Row][c.Col].ContainsCandidate(val) == false {
			continue
		}

		_ = puzzle.SetValueOnCoords(c, val)
		stats.guesses++

		switch puzzle.GetState() {
		case board.Invalid:
			_ = puzzle.SetValueOnCoords(c, board.EmptyCell)
			stats.backtracks++
		case board.Solved:
			return puzzle, nil
		default:
			solvedPuzzle, err := bt.backtrackSolve(puzzle, stats)

			if err == nil {
				return solvedPuzzle, nil
			}

			_ = puzzle.SetValueOnCoords(c, board.EmptyCell)
			stats.backtracks++
		}
	}

	return puzzle, ErrCannotSolve
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
