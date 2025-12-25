package solver

import (
	"testing"

	"github.com/DavidGudovic/sudoku/internal/core/board"
	"github.com/stretchr/testify/assert"
)

func testCases() []struct {
	name       string
	puzzle     string
	wantSolved bool
} {
	return []struct {
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
			name:       "Moderately Hard",
			puzzle:     "076000000380000000000104078000500009809307020007201030000970050000000300002000400",
			wantSolved: true,
		},
		{
			name:       "Vicious Puzzle",
			puzzle:     "097600504003000090060000000006900805700005000000030200000870003450020080000090600",
			wantSolved: true,
		},
		{
			name:       "Hard Puzzle",
			puzzle:     "200400580000208409090000000000007000309000600714000205003800006000035040000940002",
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
}

func TestBruteForceSolver(t *testing.T) {
	solver := NewBruteForceSolver()

	tests := testCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := board.FromString(tt.puzzle, false)
			boardAfterAttempt, _, err := solver.Solve(*b)

			if tt.wantSolved {
				assert.NoError(t, err)
				assert.Equal(t, board.Solved, boardAfterAttempt.State())
				return
			}

			assert.ErrorIs(t, err, ErrUnsolvablePuzzle)
			assert.Equal(t, board.Unsolved, boardAfterAttempt.State())
		})
	}
}

func TestLogicalSolver(t *testing.T) {
	solver := NewLogicalSolver()

	tests := testCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := board.FromString(tt.puzzle, false)
			boardAfterAttempt, _, err := solver.Solve(*b)

			if tt.wantSolved {
				assert.NoError(t, err)
				assert.Equal(t, board.Solved, boardAfterAttempt.State())
				return
			}

			assert.ErrorIs(t, err, ErrUnsolvablePuzzle)
			assert.Equal(t, board.Invalid, boardAfterAttempt.State())
		})
	}
}
