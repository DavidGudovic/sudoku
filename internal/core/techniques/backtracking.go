package techniques

import (
	"errors"
	"fmt"
	"time"

	"github.com/DavidGudovic/sudoku/internal/core/board"
)

const (
	impossibleCandidateCount = 10
)

var (
	ErrNoEmptyCells = errors.New("no empty cells found")
)

type backtrackStats struct {
	guesses    int
	backtracks int
}

// Backtracking algorithm attempts to solve the puzzle using brute force in nanosecond range.
// It does not produce incremental steps but returns the solved puzzle if successful.
//
// Optimizations:
//   - Minimum Remaining Values (MRV) heuristic to select the next cell to fill.
//   - Constraint propagation (handled by the board) to reduce the search space.
//   - Forward checking (also handled by the board), if any board.EmptyCell has board.NoCandidates we backtrack.
//   - Early termination upon finding a solution.
func Backtracking(puzzle *board.Board) (Step, error) {
	stats := backtrackStats{}

	if puzzle.IsUnconstrained() && puzzle.State() == board.Invalid {
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
		if !puzzle.IsUnconstrained() && puzzle.CellAt(c).ContainsCandidate(val) == false {
			continue
		}

		puzzle.MustSetValueOnCoords(c, val)
		stats.guesses++

		switch puzzle.State() {
		case board.Invalid:
			puzzle.MustSetValueOnCoords(c, board.EmptyCell)
			stats.backtracks++
		case board.Solved:
			return puzzle, nil
		default:
			solvedPuzzle, err := backtrackSolve(puzzle, stats)

			if err == nil {
				return solvedPuzzle, nil
			}

			puzzle.MustSetValueOnCoords(c, board.EmptyCell)
			stats.backtracks++
		}
	}

	return puzzle, ErrCannotSolve
}

// findSuitableCell is a helper that returns the coordinates of an empty cell with the least amount candidates
//
// In CS terms, this is the Minimum Remaining Values (MRV) heuristic
func findSuitableCell(puzzle board.Board) (board.Coordinates, error) {
	leastCandidates := impossibleCandidateCount
	var bestCell board.Coordinates

	for row := 0; row < board.Size; row++ {
		for col := 0; col < board.Size; col++ {
			c := board.MustCoordinates(row, col)

			if puzzle.CellAt(c).Value() != board.EmptyCell {
				continue
			}

			candidates := puzzle.CellAt(c).Candidates()

			if candidates.Count() == 1 {
				return c, nil
			}

			if candidates.Count() < leastCandidates {
				leastCandidates = candidates.Count()
				bestCell = c
			}
		}
	}

	if leastCandidates != impossibleCandidateCount {
		return bestCell, nil
	}

	return board.Coordinates{}, ErrNoEmptyCells
}
