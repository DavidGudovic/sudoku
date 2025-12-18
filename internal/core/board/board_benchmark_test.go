package board

import (
	"testing"
)

func BenchmarkBoard_State(t *testing.B) {
	benches := []struct {
		name  string
		board string
		want  State
	}{
		{
			name:  "Known unsolved board",
			board: "010000000007005460306400000000370001702000300000904005400701006000000500580000200",
		},
		{
			name:  "Known solved board",
			board: "637159248281347956594268173816592734429783615375614829742936581953821467168475392",
		},
	}

	for _, bm := range benches {
		t.Run(bm.name, func(t *testing.B) {
			board, _ := FromString(bm.board, false)

			t.ReportAllocs()
			t.ResetTimer()

			for i := 0; i < t.N; i++ {
				_ = board.State()
			}
		})
	}
}
