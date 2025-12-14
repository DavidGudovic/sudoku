package techniques

import (
	"testing"

	"github.com/DavidGudovic/sudoku/internal/core/board"
	"github.com/DavidGudovic/sudoku/internal/ptr"
	"github.com/stretchr/testify/assert"
)

func TestTechniques(t *testing.T) {
	tests := []struct {
		name           string
		technique      Func
		board          string
		shouldProgress bool
		expecting      Step
	}{
		{
			name:           "LastDigit (Row)",
			technique:      LastDigit,
			board:          "000000000000000000000000000000000000123406789000000000000000000000000000000000000",
			shouldProgress: true,
			expecting: Step{
				Description:   "Value 5 can only go in one place in Row 4, placing a 5 at R4C4",
				Technique:     "LastDigit (Row)",
				AffectedCells: []board.Coordinates{{Row: 4, Col: 4}},
				ReasonCells: []board.Coordinates{
					{Row: 4, Col: 0}, {Row: 4, Col: 1}, {Row: 4, Col: 2},
					{Row: 4, Col: 3}, {Row: 4, Col: 5}, {Row: 4, Col: 6},
					{Row: 4, Col: 7}, {Row: 4, Col: 8},
				},
				RemovedCandidates: 0b0000100000,
				PlacedValue:       ptr.To(5),
			},
		},
		{
			name:           "LastDigit (Column)",
			technique:      LastDigit,
			board:          "000600000000500000000700000000000000000100000000300000000900000000200000000800000",
			shouldProgress: true,
			expecting: Step{
				Description:   "Value 4 can only go in one place in Col 3, placing a 4 at R3C3",
				Technique:     "LastDigit (Column)",
				AffectedCells: []board.Coordinates{{Row: 3, Col: 3}},
				ReasonCells: []board.Coordinates{
					{Row: 0, Col: 3}, {Row: 1, Col: 3}, {Row: 2, Col: 3},
					{Row: 4, Col: 3}, {Row: 5, Col: 3}, {Row: 6, Col: 3},
					{Row: 7, Col: 3}, {Row: 8, Col: 3},
				},
				RemovedCandidates: 0b0000010000,
				PlacedValue:       ptr.To(4),
			},
		},
		{
			name:           "LastDigit (Box)",
			technique:      LastDigit,
			board:          "000000000000000000000000000000123000000604000000789000000000000000000000000000000",
			shouldProgress: true,
			expecting: Step{
				Description:   "Value 5 can only go in one place in Box 4, placing a 5 at R4C4",
				Technique:     "LastDigit (Box)",
				AffectedCells: []board.Coordinates{{Row: 4, Col: 4}},
				ReasonCells: []board.Coordinates{
					{Row: 3, Col: 3}, {Row: 3, Col: 4}, {Row: 3, Col: 5},
					{Row: 4, Col: 3}, {Row: 4, Col: 5},
					{Row: 5, Col: 3}, {Row: 5, Col: 4}, {Row: 5, Col: 5},
				},
				RemovedCandidates: 0b0000100000,
				PlacedValue:       ptr.To(5),
			},
		},
		{
			name:           "LastDigit (No Progress)",
			technique:      LastDigit,
			board:          "530070000600195000098000060800060003400803001700020006060000280000419005000080079",
			shouldProgress: false,
			expecting:      Step{},
		},
		{
			name:           "NakedSingle (Progress)",
			technique:      NakedSingle,
			board:          "002000000006000000007000000008000000004000000003000000000000150000000000000000000",
			shouldProgress: true,
			expecting: Step{
				Description:       "The candidate 9 is the only one left at R6C2, placing a 9",
				Technique:         "NakedSingle",
				AffectedCells:     []board.Coordinates{{Row: 6, Col: 2}},
				ReasonCells:       allPeers(board.Coordinates{Row: 6, Col: 2}).Excluding(board.Coordinates{Row: 6, Col: 2}),
				RemovedCandidates: 0b1000000000,
				PlacedValue:       ptr.To(9),
			},
		},
		{
			name:           "NakedSingle (No Progress)",
			technique:      NakedSingle,
			board:          "690583010105090803830010500063870100058421036210630008526947381389152647001368000",
			shouldProgress: false,
			expecting:      Step{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := board.FromString(tt.board, false)
			step, err := tt.technique.Apply(b)

			if tt.shouldProgress == false {
				assert.False(t, step.MadeProgress())
				assert.Equal(t, tt.expecting, step)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expecting, step)
			assert.True(t, step.MadeProgress())
		})
	}
}
