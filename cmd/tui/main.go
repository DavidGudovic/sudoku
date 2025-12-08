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

	boardFromAString.SetValue(1, 6, board.Nine)
	boardFromAString.SetValue(5, 2, board.One)
	boardFromAString.SetValue(8, 5, board.Two)
	boardFromAString.SetValue(6, 7, board.Five)
	boardFromAString.SetValue(2, 7, board.Six)
	boardFromAString.SetValue(3, 1, board.Four)

	fmt.Println(boardFromAString.ToString(true))
}
