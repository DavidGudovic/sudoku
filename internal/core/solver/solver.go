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

func NewSudokuSolver(techniques []techniques.Technique) *SudokuSolver {
	return &SudokuSolver{
		techniques: techniques,
	}
}

func NewBruteForceSolver() *SudokuSolver {
	return &SudokuSolver{
		techniques: []techniques.Technique{
			techniques.Backtracking{},
		},
	}
}

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
					s.PropagateCandidates(&puzzle)
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

func (s *SudokuSolver) PropagateCandidates(_ *board.Board) {
	// TODO: Implement
}
