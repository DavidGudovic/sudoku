package solver

import (
	"testing"

	"github.com/DavidGudovic/sudoku/internal/core/board"
)

func BenchmarkBruteForceSolver(b *testing.B) {
	solver := NewBruteForceSolver()

	benchmarks := []struct {
		name   string
		puzzle string
	}{
		{
			name:   "Easy",
			puzzle: "081000670000007050003280000030000890708301260002800104010530040350000000890004000",
		},
		{
			name:   "Vicious",
			puzzle: "097600504003000090060000000006900805700005000000030200000870003450020080000090600",
		},
		{
			name:   "Hardest",
			puzzle: "206050470070000002300000000000180000400700905000000810903070600000005030160200009",
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			original, _ := board.FromString(bm.puzzle)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				puzzle := *original
				_, _, _ = solver.Solve(puzzle)
			}
		})
	}
}
