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
		technique      FuncAdapter
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
				Description:       "Value 5 can only go in one place in Row 4, placing a 5 at R4C4",
				Technique:         "LastDigit (Row)",
				AffectedCells:     []board.Coordinates{{Row: 4, Col: 4}},
				ReasonCells:       Peers.Of(board.Coordinates{Row: 4, Col: 4}).Across(Row).Slice(),
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
				Description:       "Value 4 can only go in one place in Col 3, placing a 4 at R3C3",
				Technique:         "LastDigit (Column)",
				AffectedCells:     []board.Coordinates{{Row: 3, Col: 3}},
				ReasonCells:       Peers.Of(board.Coordinates{Row: 3, Col: 3}).Across(Column).Slice(),
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
				Description:       "Value 5 can only go in one place in Box 4, placing a 5 at R4C4",
				Technique:         "LastDigit (Box)",
				AffectedCells:     []board.Coordinates{{Row: 4, Col: 4}},
				ReasonCells:       Peers.Of(board.Coordinates{Row: 4, Col: 4}).Across(Box).Slice(),
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
				ReasonCells:       Peers.Of(board.Coordinates{Row: 6, Col: 2}).Across(AllScopes...).Slice(),
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
		{
			name:           "NakedPair (Progress)",
			technique:      NakedPair,
			board:          "000524768784916235265700491047000006020807143000000927010070309070600014000301672",
			shouldProgress: true,
			expecting: Step{
				Description:       "NakedPair found at R3C6,R3C7, removing candidates 5,8 from mutual peers",
				Technique:         "NakedPair",
				ReasonCells:       []board.Coordinates{{Row: 3, Col: 6}, {Row: 3, Col: 7}},
				AffectedCells:     []board.Coordinates{{Row: 3, Col: 0}, {Row: 3, Col: 4}, {Row: 3, Col: 5}},
				RemovedCandidates: 0b0100100000,
			},
		},
		{
			name:           "NakedPair (No Progress)",
			technique:      NakedPair,
			board:          "975421386148563700632879154006200430004300000390004000000940563409130278003002941",
			shouldProgress: false,
			expecting:      Step{},
		},
		{
			name:           "NakedTriple",
			technique:      NakedTriple,
			board:          "000000060431628795500001020642319587715286934003754216004000072007002350000800040",
			shouldProgress: true,
			expecting: Step{
				Description:       "NakedTriple found at R8C2,R8C6,R8C8, removing candidates 1,6,9 from mutual peers",
				Technique:         "NakedTriple",
				ReasonCells:       []board.Coordinates{{Row: 8, Col: 2}, {Row: 8, Col: 6}, {Row: 8, Col: 8}},
				AffectedCells:     []board.Coordinates{{Row: 8, Col: 0}, {Row: 8, Col: 1}, {Row: 8, Col: 4}},
				RemovedCandidates: 0b1001000010,
			},
		},
		{
			name:           "NakedTriple (No Progress)",
			technique:      NakedTriple,
			board:          "000000000000000000000000000000000000123406789000000000000000000000000000000000000",
			shouldProgress: false,
			expecting:      Step{},
		},
		{
			name:           "NakedQuad",
			technique:      NakedQuad,
			board:          "381962700465387921927541600010473000040608000038209000276134009893725416154896372",
			shouldProgress: true,
			expecting: Step{
				Description:       "NakedQuad found at R3C6,R3C8,R4C6,R5C6, removing candidates 1,2,5,8 from mutual peers",
				Technique:         "NakedQuad",
				ReasonCells:       []board.Coordinates{{Row: 3, Col: 6}, {Row: 3, Col: 8}, {Row: 4, Col: 6}, {Row: 5, Col: 6}},
				AffectedCells:     []board.Coordinates{{Row: 3, Col: 7}, {Row: 4, Col: 7}, {Row: 4, Col: 8}, {Row: 5, Col: 7}, {Row: 5, Col: 8}},
				RemovedCandidates: 0b0100100110,
			},
		},
		{
			name:           "NakedQuad (No Progress)",
			technique:      NakedQuad,
			board:          "000000000000000000000000000000000000123406789000000000000000000000000000000000000",
			shouldProgress: false,
			expecting:      Step{},
		},
		{
			name:           "HiddenSingle (Row)",
			technique:      HiddenSingle,
			board:          "000000020000000000857034000000000000000000000000000000000000000000000000000000000",
			shouldProgress: true,
			expecting: Step{
				Description:       "None of the empty cells in Row 2 can hold a 2 except R2C3, placing a 2",
				Technique:         "HiddenSingle (Row)",
				ReasonCells:       []board.Coordinates{{Row: 2, Col: 6}, {Row: 2, Col: 7}, {Row: 2, Col: 8}},
				AffectedCells:     []board.Coordinates{{Row: 2, Col: 3}},
				RemovedCandidates: 0b0000000100,
				PlacedValue:       ptr.To(2),
			},
		},
		{
			name:           "HiddenSingle (Box)",
			technique:      HiddenSingle,
			board:          "000000000000000000000000000000000000000000000000000000000000072000000431060000500",
			shouldProgress: true,
			expecting: Step{
				Description:       "None of the empty cells in Box 8 can hold a 6 except R6C6, placing a 6",
				Technique:         "HiddenSingle (Box)",
				ReasonCells:       []board.Coordinates{{Row: 8, Col: 7}, {Row: 8, Col: 8}},
				AffectedCells:     []board.Coordinates{{Row: 6, Col: 6}},
				RemovedCandidates: 0b0001000000,
				PlacedValue:       ptr.To(6),
			},
		},
		{
			name:           "HiddenSingle (No Progress)",
			technique:      HiddenSingle,
			board:          "975421386148563700632879154006200430004300000390004000000940563409130278003002941",
			shouldProgress: false,
			expecting:      Step{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := board.FromString(tt.board, false)
			step, err := tt.technique(b)

			if tt.shouldProgress == false {
				assert.Error(t, err)
				assert.Equal(t, ErrCannotProgress, err)
				assert.False(t, step.MadeProgress())
				assert.Equal(t, tt.expecting, step)
				return
			}

			// Should have made progress with no errors and an appropriate step
			assert.NoError(t, err)
			assert.Equal(t, tt.expecting, step)
			assert.True(t, step.MadeProgress())

			// Should have affected the expected cells with correct placements/removals
			for _, affected := range step.AffectedCells {
				cell := b.CellAt(affected)
				if step.PlacedValue != nil {
					assert.Equal(t, *step.PlacedValue, cell.Value())
				}

				if step.RemovedCandidates != board.NoCandidates {
					IntersectingCandidates := cell.Candidates() & step.RemovedCandidates
					assert.Empty(t, IntersectingCandidates)
				}
			}
		})
	}
}
