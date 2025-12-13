package solver

import (
	"errors"

	"github.com/DavidGudovic/sudoku/internal/core/board"
	"github.com/DavidGudovic/sudoku/internal/core/solver/techniques"
)

var (
	ErrUnsolvablePuzzle = errors.New("puzzle is unsolvable with the given techniques")
)

type Solver interface {
	Solve(puzzle board.Board) (board.Board, []techniques.Step, error)
}

type SudokuSolver struct {
	techniques []techniques.Technique
}

// NewSudokuSolver creates a SudokuSolver with the specified techniques.
func NewSudokuSolver(techniques []techniques.Technique) *SudokuSolver {
	return &SudokuSolver{
		techniques: techniques,
	}
}

// NewBruteForceSolver creates a SudokuSolver that uses only the backtracking technique.
func NewBruteForceSolver() *SudokuSolver {
	return &SudokuSolver{
		techniques: []techniques.Technique{
			techniques.Backtracking{},
		},
	}
}

// NewLogicalSolver creates a SudokuSolver that uses a set of human-like logical techniques.
func NewLogicalSolver() *SudokuSolver {
	return &SudokuSolver{
		techniques: []techniques.Technique{
			techniques.NakedSingle{},
			techniques.HiddenSingle{},
			techniques.NakedPair{},
			techniques.HiddenPair{},
			techniques.PointingPair{},
			techniques.XWing{},
			techniques.Skyscraper{},
		},
	}
}

// Solve attempts to solve the given Sudoku puzzle using the configured techniques.
// It returns the solved board, a list of steps taken to solve it, or an error if unsolvable.
func (s *SudokuSolver) Solve(puzzle board.Board) (board.Board, []techniques.Step, error) {
	var steps []techniques.Step

	for {
		progressMade := false

		// Attempt each technique in order
		for _, technique := range s.techniques {
			step, err := technique.Apply(&puzzle)

			if err != nil {
				return puzzle, nil, err
			}

			if step.MadeProgress() {
				steps = append(steps, step)
				progressMade = true // If any technique made progress, the puzzle is still solvable

				if step.PlacedValue != nil {
					break // Previous (cheaper) techniques may now be applicable
				}
			}
		}

		// A full cycle completed with no progress
		// means further attempts are futile, or the puzzle is solved
		if !progressMade {
			break
		}
	}

	if puzzle.GetState() == board.Solved {
		return puzzle, steps, nil
	}

	return puzzle, nil, ErrUnsolvablePuzzle
}
