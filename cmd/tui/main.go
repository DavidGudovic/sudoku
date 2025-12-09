package main

import (
	"fmt"

	"github.com/DavidGudovic/sudoku/internal/core/board"
	"github.com/DavidGudovic/sudoku/internal/core/solver"
)

func main() {
	emptyBoard := board.NewBoard()

	boardStringRepresentation := emptyBoard.ToString(true)

	boardFromAString, err := board.FromString(boardStringRepresentation)

	if err != nil {
		panic(err)
	}

	fmt.Println(boardFromAString.ToString(true))

	boardFromAString.SetValueOnCoords(1, 6, board.Nine)
	boardFromAString.SetValueOnCoords(5, 2, board.One)
	boardFromAString.SetValueOnCoords(8, 5, board.Two)
	boardFromAString.SetValueOnCoords(6, 7, board.Five)
	boardFromAString.SetValueOnCoords(2, 7, board.Six)
	boardFromAString.SetValueOnCoords(3, 1, board.Four)

	fmt.Println(boardFromAString.ToString(true))

	bfv := solver.NewBacktrackingSolver(solver.NewValidator())

	state, err := bfv.Validator.ValidateBoardState(boardFromAString)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Board is solved:", state)

	solvedBoardString := "637159248281347956594268173816592734429783615375614829742936581953821467168475392"

	solvedBoard, err := board.FromString(solvedBoardString)

	if err != nil {
		panic(err)
	}

	boardState, err := bfv.Validator.ValidateBoardState(solvedBoard)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("The known solved Board is solved:", boardState)

	invalidBoard := board.NewBoard()
	invalidBoard.SetValueOnCoords(1, 6, board.Nine)
	invalidBoard.SetValueOnCoords(2, 6, board.Nine)
	invalidBoard.SetValueOnCoords(2, 2, board.Six)
	invalidBoard.SetValueOnCoords(2, 3, board.Six)
	invalidBoard.SetValueOnCoords(5, 2, board.One)
	invalidBoard.SetValueOnCoords(6, 1, board.Two)
	invalidBoard.SetValueOnCoords(7, 4, board.Three)

	boardState, err = bfv.Validator.ValidateBoardState(invalidBoard)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("The invalid Board is solved:", boardState)

}
