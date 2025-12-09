package main

import (
	"fmt"

	"github.com/DavidGudovic/sudoku/internal/core/board"
)

func main() {
	emptyBoard := board.NewBoard()

	boardStringRepresentation := emptyBoard.ToString(true)

	boardFromAString, err := board.FromString(boardStringRepresentation)

	if err != nil {
		panic(err)
	}

	fmt.Println(boardFromAString.ToString(true))

	_ = boardFromAString.SetValueOnCoords(1, 6, 9)
	_ = boardFromAString.SetValueOnCoords(5, 2, 1)
	_ = boardFromAString.SetValueOnCoords(8, 5, 2)
	_ = boardFromAString.SetValueOnCoords(6, 7, 5)
	_ = boardFromAString.SetValueOnCoords(2, 7, 6)
	_ = boardFromAString.SetValueOnCoords(3, 1, 4)

	fmt.Println(boardFromAString.ToString(true))

	state := boardFromAString.ValidateState()

	fmt.Println("State of known unsolved grid:", state)

	solvedBoardString := "637159248281347956594268173816592734429783615375614829742936581953821467168475392"

	solvedBoard, err := board.FromString(solvedBoardString)

	if err != nil {
		panic(err)
	}

	boardState := solvedBoard.ValidateState()

	fmt.Println("State of known solved grid:", boardState)

	invalidBoard := board.NewBoard()
	_ = invalidBoard.SetValueOnCoords(1, 6, 9)
	_ = invalidBoard.SetValueOnCoords(2, 6, 9)
	_ = invalidBoard.SetValueOnCoords(2, 2, 6)
	_ = invalidBoard.SetValueOnCoords(2, 3, 6)
	_ = invalidBoard.SetValueOnCoords(5, 2, 3)
	_ = invalidBoard.SetValueOnCoords(6, 1, 3)
	_ = invalidBoard.SetValueOnCoords(7, 4, 5)

	boardState = invalidBoard.ValidateState()

	fmt.Println("State of known invalid grid:", boardState)

}
