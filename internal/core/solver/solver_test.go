package solver

import (
	"testing"

	"github.com/DavidGudovic/sudoku/internal/core/board"
	"github.com/stretchr/testify/assert"
)

func TestBruteForceSolver(t *testing.T) {
	solver := NewBruteForceSolver()

	tests := []struct {
		name       string
		puzzle     string
		wantSolved bool
	}{
		{
			name:       "Moderately easy",
			puzzle:     "081000670000007050003280000030000890708301260002800104010530040350000000890004000",
			wantSolved: true,
		},
		{
			name:       "Vicious Puzzle",
			puzzle:     "097600504003000090060000000006900805700005000000030200000870003450020080000090600",
			wantSolved: true,
		},
		{
			name:       "Hardest Puzzle",
			puzzle:     "206050470070000002300000000000180000400700905000000810903070600000005030160200009",
			wantSolved: true,
		},
		{
			name:       "Unsolvable Puzzle",
			puzzle:     "100900500000610080020000010600002800030000006000160020810000000000028061050001008",
			wantSolved: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := board.FromString(tt.puzzle)
			boardAfterAttempt, _, err := solver.Solve(*b)

			if tt.wantSolved {
				assert.NoError(t, err)
				assert.Equal(t, board.Solved, boardAfterAttempt.GetState())
				return
			}

			assert.Error(t, err, ErrUnsolvablePuzzle)
			assert.Equal(t, board.Unsolved, boardAfterAttempt.GetState())
		})
	}
}
