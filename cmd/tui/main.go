package main

import (
	"fmt"

	"github.com/DavidGudovic/sudoku/internal/core/board"
	"github.com/DavidGudovic/sudoku/internal/core/solver"
)

func main() {
	puzzle, _ := board.FromString("206050470070000002300000000000180000400700905000000810903070600000005030160200009")

	s := solver.NewBruteForceSolver()

	solvedBoard, steps, err := s.Solve(*puzzle)

	if err != nil {
		panic(err)
	}

	fmt.Print("Given puzzle:\n")
	fmt.Println(puzzle.ToString(false))
	fmt.Print("\nSolved puzzle:\n")
	fmt.Println(solvedBoard.ToString(false))
	fmt.Print("\nSolving steps:\n")
	fmt.Println(steps[0].Description)
}
