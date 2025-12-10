package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoard_GetState(t *testing.T) {
	tests := []struct {
		name  string
		board string
		want  State
	}{
		{
			name:  "Known unsolved valid board",
			board: "63715920_2_3_482813479565902681738165927344290_1_2_3_4_5_6_7_8_983615375614829742936581053821467168475392",
			want:  Unsolved,
		},
		{
			name:  "Known solved board",
			board: "637159248281347956594268173816592734429783615375614829742936581953821467168475392",
			want:  Solved,
		},
		{
			name:  "Known invalid board",
			board: "637159228281347956894268173816592334429783615375614829742936581453824467168675992",
			want:  Invalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board, _ := FromString(tt.board)
			assert.Equal(t, tt.want, board.GetState())
		})
	}
}
