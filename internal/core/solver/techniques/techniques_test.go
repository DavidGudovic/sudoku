package techniques

import (
	"testing"

	"github.com/DavidGudovic/sudoku/internal/core/board"
	"github.com/DavidGudovic/sudoku/internal/ptr"
	"github.com/stretchr/testify/assert"
)

func TestTechniques(t *testing.T) {
	tests := []struct {
		name      string
		technique Technique
		expecting Step
		board     string
	}{
		{
			name:      "LastDigit",
			technique: LastDigit{},
			board:     "000000000000000000000000000000000000123406789000000000000000000000000000000000000",
			expecting: Step{
				Description:   "Last Digit Technique Applied in Row 4, placing a 5 at R4C4.",
				Technique:     "LastDigit",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := board.FromString(tt.board, false)
			step, err := tt.technique.Apply(b)

			assert.NoError(t, err)
			assert.Equal(t, tt.expecting, step)
		})
	}
}
