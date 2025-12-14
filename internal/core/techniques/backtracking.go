package techniques

import (
	"errors"
	"fmt"
	"time"

	"github.com/DavidGudovic/sudoku/internal/core/board"
)

type backtrackStats struct {
	guesses    int
	backtracks int
}

// Backtracking algorithm attempts to solve the puzzle using brute force.
// It does not produce incremental steps but returns the solved puzzle if successful.
func Backtracking(puzzle *board.Board) (Step, error) {
	stats := backtrackStats{}

	if puzzle.IsFaux() && puzzle.GetState() == board.Invalid {
		return Step{}, ErrCannotSolve
	}

	start := time.Now()
	solvedPuzzle, err := backtrackSolve(*puzzle, &stats)
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
func backtrackSolve(puzzle board.Board, stats *backtrackStats) (board.Board, error) {
	c, err := findSuitableCell(puzzle)

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
			solvedPuzzle, err := backtrackSolve(puzzle, stats)

			if err == nil {
				return solvedPuzzle, nil
			}

			_ = puzzle.SetValueOnCoords(c, board.EmptyCell)
			stats.backtracks++
		}
	}

	return puzzle, ErrCannotSolve
}

// findSuitableCell is a helper that returns the coordinates of an empty cell with the least amount candidates
// Searches left to right, top to bottom.
func findSuitableCell(puzzle board.Board) (board.Coordinates, error) {
	leastCandidates := 10 // Impossible value.
	var bestCell board.Coordinates

	for row := 0; row < board.Size; row++ {
		for col := 0; col < board.Size; col++ {
			if puzzle.Cells[row][col].Value() != board.EmptyCell {
				continue
			}

			candidates := puzzle.CellAt(board.Coordinates{Row: row, Col: col}).Candidates()

			if candidates.Count() == 1 {
				return board.Coordinates{Row: row, Col: col}, nil
			}

			if candidates.Count() < leastCandidates {
				leastCandidates = candidates.Count()
				bestCell = board.Coordinates{Row: row, Col: col}
			}
		}
	}

	if leastCandidates != 10 {
		return bestCell, nil
	}

	return board.Coordinates{}, errors.New("no empty cells found")
}
