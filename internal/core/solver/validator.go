package solver

import (
	"errors"

	"github.com/DavidGudovic/sudoku/internal/core/board"
)

var (
	ErrInvalidBoardState = errors.New("invalid board state")
)

type BruteForceValidator struct{}

func NewValidator() *BruteForceValidator {
	return &BruteForceValidator{}
}

func (v *BruteForceValidator) ValidateBoardState(b *board.Board) (bool, error) {
	allRowsSolved, err := validateRows(b)
	if err != nil {
		return false, err
	}

	allColsSolved, err := validateCols(b)
	if err != nil {
		return false, err
	}

	allBoxesSolved, err := validateBoxes(b)
	if err != nil {
		return false, err
	}

	solved := allRowsSolved && allColsSolved && allBoxesSolved
	return solved, nil
}

func validateRows(b *board.Board) (bool, error) {
	var seen [board.MaxValue + 1]bool
	allRowsSolved := true

	for _, rowArray := range b.Cells {
		for _, cell := range rowArray {
			if cell.Value == board.EmptyCell {
				continue
			}

			if seen[cell.Value] {
				return false, ErrInvalidBoardState
			}

			seen[cell.Value] = true
		}
		allRowsSolved = allRowsSolved && allValuesSeen(seen)
		seen = [board.MaxValue + 1]bool{}
	}

	return allRowsSolved, nil
}

func validateCols(b *board.Board) (bool, error) {
	var seen [board.MaxValue + 1]bool
	allColsSolved := true

	for col := 0; col < board.ColCount; col++ {
		for row := 0; row < board.RowCount; row++ {
			value := b.Cells[row][col].Value

			if value == board.EmptyCell {
				continue
			}

			if seen[value] {
				return false, ErrInvalidBoardState
			}

			seen[value] = true
		}

		allColsSolved = allColsSolved && allValuesSeen(seen)
		seen = [board.MaxValue + 1]bool{}
	}

	return allColsSolved, nil
}

func validateBoxes(b *board.Board) (bool, error) {
	var seen [board.MaxValue + 1]bool
	allBoxesSolved := true

	for boxRow := range board.BoxRowCount {
		for boxCol := range board.BoxColCount {

			for rowInBox := 0; rowInBox < board.BoxRowCount; rowInBox++ {
				for colInBox := 0; colInBox < board.BoxColCount; colInBox++ {
					rowInSudoku := boxRow*board.BoxRowCount + rowInBox
					colInSudoku := boxCol*board.BoxColCount + colInBox

					value := b.Cells[rowInSudoku][colInSudoku].Value

					if value == board.EmptyCell {
						continue
					}

					if seen[value] {
						return false, ErrInvalidBoardState
					}

					seen[value] = true
				}
			}

			allBoxesSolved = allBoxesSolved && allValuesSeen(seen)
			seen = [board.MaxValue + 1]bool{}
		}
	}

	return allBoxesSolved, nil
}

func allValuesSeen(seen [board.MaxValue + 1]bool) bool {
	for _, value := range board.AllCellValues {
		if !seen[value] {
			return false
		}
	}

	return true
}
