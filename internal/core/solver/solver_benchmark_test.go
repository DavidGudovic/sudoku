package solver

import (
	"testing"

	"github.com/DavidGudovic/sudoku/internal/core/board"
	"github.com/DavidGudovic/sudoku/internal/core/techniques"
)

type Benchmark struct {
	name   string
	puzzle string
}

func Setup() []Benchmark {
	return []Benchmark{
		{
			name:   "Easy",
			puzzle: "081000670000007050003280000030000890708301260002800104010530040350000000890004000",
		},
		{
			name:   "Hard",
			puzzle: "200400580000208409090000000000007000309000600714000205003800006000035040000940002",
		},
		{
			name:   "Vicious",
			puzzle: "097600504003000090060000000006900805700005000000030200000870003450020080000090600",
		},
		{
			name:   "Beyond Hell",
			puzzle: "206050470070000002300000000000180000400700905000000810903070600000005030160200009",
		},
		{
			name:   "AlEscargot (2006)",
			puzzle: "100007090030020008009600500005300900010080002600004000300000010040000007007000300",
		},
	}
}

// BenchmarkBruteForceSolver benchmarks the brute force solver with all the overhead of Solver orchestration, Step generation, and similar.
func BenchmarkBruteForceSolver(b *testing.B) {
	solver := NewBruteForceSolver()
	benchmarks := Setup()

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			original, _ := board.FromString(bm.puzzle, false)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _, _ = solver.Solve(*original)
			}
		})
	}
}

// BenchmarkBacktracking benchmarks the backtracking algorithm used in a standard BruteForceSolver, eliminating Solver orchestration but keeping Step generation and safety checks overhead
func BenchmarkBacktracking(b *testing.B) {
	benchmarks := Setup()

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			original, _ := board.FromString(bm.puzzle, false)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				temp := *original
				_, _ = techniques.Backtracking(&temp)
			}
		})
	}
}

func BenchmarkLogicalSolving(b *testing.B) {
	solver := NewLogicalSolver()
	benchmarks := Setup()

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			original, _ := board.FromString(bm.puzzle, false)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _, err := solver.Solve(*original)

				if err != nil {
					b.Fatal("Logical solver failed to solve the puzzle:", err)
				}
			}
		})
	}
}
