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
			board: "63715920*2*3*482813479565902681738165927344290*1*2*3*4*5*6*7*8*983615375614829742936581053821467168475392",
			want:  Unsolved,
		},
		{
			name:  "Known unsolved valid board 2",
			board: "006003020070004000100006970002008090700030061030600000408000005000000002001040730",
			want:  Unsolved,
		},
		{
			name:  "Known unsolved valid board 3",
			board: "010000000007005460306400000000370001702000300000904005400701006000000500580000200",
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
			isFauxBoard := tt.want == Invalid
			board, _ := FromString(tt.board, isFauxBoard)
			assert.Equal(t, tt.want, board.GetState())
		})
	}
}

func TestCoordsFromIndex(t *testing.T) {
	tests := []struct {
		name    string
		index   int
		wantRow int
		wantCol int
		wantErr bool
	}{
		{
			name:    "Valid Index",
			index:   54,
			wantRow: 6,
			wantCol: 0,
			wantErr: false,
		},
		{
			name:    "Invalid Index",
			index:   -1,
			wantRow: 0,
			wantCol: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := CoordsFromIndex(tt.index)

			if tt.wantErr {
				assert.Error(t, err, ErrIndexOutOfBounds)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantRow, c.Row)
			assert.Equal(t, tt.wantCol, c.Col)
		})
	}
}

func TestBoard_GetSetValue(t *testing.T) {
	tests := []struct {
		name    string
		board   *Board
		index   int
		wantErr bool
		value   int
	}{
		{
			name:    "Set Value",
			board:   NewBoard(),
			index:   54,
			wantErr: false,
			value:   1,
		},
		{
			name:    "Set Illegal Value",
			board:   NewBoard(),
			index:   54,
			wantErr: true,
			value:   10,
		},
		{
			name:    "Set Illegal Index",
			board:   NewBoard(),
			index:   -1,
			wantErr: true,
			value:   5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.board.SetValueOnIndex(tt.index, tt.value)

			if tt.wantErr {
				assert.Error(t, err, ErrIndexOutOfBounds)
				return
			}

			assert.NoError(t, err)
			value, _ := tt.board.GetValueByIndex(tt.index)
			assert.Equal(t, tt.value, value)
		})
	}
}

func TestBoard_SerializeSerializer(t *testing.T) {
	board := NewBoard()
	newBoard, err := FromString(board.ToString(true), false)

	assert.NoError(t, err)
	assert.Equal(t, newBoard, board)
}

func TestBoard_IsFaux(t *testing.T) {
	tests := []struct {
		name     string
		board    *Board
		wantFaux bool
	}{
		{
			name:     "Normal board",
			board:    NewBoard(),
			wantFaux: false,
		},
		{
			name:     "Faux board",
			board:    NewFauxBoard(),
			wantFaux: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantFaux, tt.board.IsFaux())
		})
	}
}

func TestBoard_ConstraintEnforcement(t *testing.T) {
	tests := []struct {
		name                 string
		board                *Board
		setValue             struct{ row, col, value int }
		checkCell            struct{ row, col int }
		wantCandidateRemoved bool
	}{
		{
			name:                 "Normal board enforces constraints",
			board:                NewBoard(),
			setValue:             struct{ row, col, value int }{0, 0, 5},
			checkCell:            struct{ row, col int }{0, 1},
			wantCandidateRemoved: true,
		},
		{
			name:                 "Faux board does not enforce constraints",
			board:                NewFauxBoard(),
			setValue:             struct{ row, col, value int }{0, 0, 5},
			checkCell:            struct{ row, col int }{0, 1},
			wantCandidateRemoved: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := NewCoordinates(tt.setValue.row, tt.setValue.col)
			err := tt.board.SetValueOnCoords(c, tt.setValue.value)
			assert.NoError(t, err)

			checkCell := tt.board.Cells[tt.checkCell.row][tt.checkCell.col]
			hasCandidate := checkCell.ContainsCandidate(tt.setValue.value)

			if tt.wantCandidateRemoved {
				assert.False(t, hasCandidate)
			} else {
				assert.True(t, hasCandidate)
			}
		})
	}
}

func TestFromString_FauxBoard(t *testing.T) {
	tests := []struct {
		name        string
		boardString string
		isFauxBoard bool
		wantFaux    bool
	}{
		{
			name:        "Normal board from string",
			boardString: "637159248281347956594268173816592734429783615375614829742936581953821467168475392",
			isFauxBoard: false,
			wantFaux:    false,
		},
		{
			name:        "Faux board from string",
			boardString: "637159248281347956594268173816592734429783615375614829742936581953821467168475392",
			isFauxBoard: true,
			wantFaux:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board, err := FromString(tt.boardString, tt.isFauxBoard)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantFaux, board.IsFaux())
		})
	}
}
