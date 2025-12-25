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
		technique      Technique
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
				Technique:         "LastDigit (Row)",
				AffectedCells:     NoPeers.Including([]board.Coordinates{{Row: 4, Col: 4}}...),
				ReasonCells:       Peers.Of(board.Coordinates{Row: 4, Col: 4}).Across(Row),
				RemovedCandidates: board.MustCandidateSet(5),
				PlacedValue:       ptr.To(5),
			},
		},
		{
			name:           "LastDigit (Column)",
			technique:      LastDigit,
			board:          "000600000000500000000700000000000000000100000000300000000900000000200000000800000",
			shouldProgress: true,
			expecting: Step{
				Technique:         "LastDigit (Column)",
				AffectedCells:     NoPeers.Including([]board.Coordinates{{Row: 3, Col: 3}}...),
				ReasonCells:       Peers.Of(board.Coordinates{Row: 3, Col: 3}).Across(Column),
				RemovedCandidates: board.MustCandidateSet(4),
				PlacedValue:       ptr.To(4),
			},
		},
		{
			name:           "LastDigit (Box)",
			technique:      LastDigit,
			board:          "000000000000000000000000000000123000000604000000789000000000000000000000000000000",
			shouldProgress: true,
			expecting: Step{
				Technique:         "LastDigit (Box)",
				AffectedCells:     NoPeers.Including([]board.Coordinates{{Row: 4, Col: 4}}...),
				ReasonCells:       Peers.Of(board.Coordinates{Row: 4, Col: 4}).Across(Box),
				RemovedCandidates: board.MustCandidateSet(5),
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
				Technique:         "NakedSingle",
				AffectedCells:     NoPeers.Including([]board.Coordinates{{Row: 6, Col: 2}}...),
				ReasonCells:       Peers.Of(board.Coordinates{Row: 6, Col: 2}).Across(AllScopes...),
				RemovedCandidates: board.MustCandidateSet(9),
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
				Technique:         "NakedPair",
				ReasonCells:       NoPeers.Including([]board.Coordinates{{Row: 3, Col: 6}, {Row: 3, Col: 7}}...),
				AffectedCells:     NoPeers.Including([]board.Coordinates{{Row: 3, Col: 0}, {Row: 3, Col: 4}, {Row: 3, Col: 5}}...),
				RemovedCandidates: board.MustCandidateSet(5, 8),
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
				Technique:         "NakedTriple",
				ReasonCells:       NoPeers.Including([]board.Coordinates{{Row: 8, Col: 2}, {Row: 8, Col: 6}, {Row: 8, Col: 8}}...),
				AffectedCells:     NoPeers.Including([]board.Coordinates{{Row: 8, Col: 0}, {Row: 8, Col: 1}, {Row: 8, Col: 4}}...),
				RemovedCandidates: board.MustCandidateSet(1, 6, 9),
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
				Technique:         "NakedQuad",
				ReasonCells:       NoPeers.Including([]board.Coordinates{{Row: 3, Col: 6}, {Row: 3, Col: 8}, {Row: 4, Col: 6}, {Row: 5, Col: 6}}...),
				AffectedCells:     NoPeers.Including([]board.Coordinates{{Row: 3, Col: 7}, {Row: 4, Col: 7}, {Row: 4, Col: 8}, {Row: 5, Col: 7}, {Row: 5, Col: 8}}...),
				RemovedCandidates: board.MustCandidateSet(1, 2, 5, 8),
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
				Technique:         "HiddenSingle (Row)",
				ReasonCells:       NoPeers.Including([]board.Coordinates{{Row: 2, Col: 6}, {Row: 2, Col: 7}, {Row: 2, Col: 8}}...),
				AffectedCells:     NoPeers.Including([]board.Coordinates{{Row: 2, Col: 3}}...),
				RemovedCandidates: board.MustCandidateSet(2),
				PlacedValue:       ptr.To(2),
			},
		},
		{
			name:           "HiddenSingle (Box)",
			technique:      HiddenSingle,
			board:          "000000000000000000000000000000000000000000000000000000000000072000000431060000500",
			shouldProgress: true,
			expecting: Step{
				Technique:         "HiddenSingle (Box)",
				ReasonCells:       NoPeers.Including([]board.Coordinates{{Row: 8, Col: 7}, {Row: 8, Col: 8}}...),
				AffectedCells:     NoPeers.Including([]board.Coordinates{{Row: 6, Col: 6}}...),
				RemovedCandidates: board.MustCandidateSet(6),
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
		{
			name:           "HiddenPair",
			technique:      HiddenPair,
			board:          "000000020504830000000000700000000000000000000000000000000000000000000000000000000",
			shouldProgress: true,
			expecting: Step{
				Technique:         "HiddenPair",
				ReasonCells:       NoPeers.Including([]board.Coordinates{{Row: 1, Col: 1}, {Row: 1, Col: 5}}...),
				AffectedCells:     NoPeers.Including([]board.Coordinates{{Row: 1, Col: 1}, {Row: 1, Col: 5}}...),
				RemovedCandidates: board.MustCandidateSet(1, 6, 9),
			},
		},
		{
			name:           "HiddenPair (No Progress)",
			technique:      HiddenPair,
			board:          "975421386148563700632879154006200430004300000390004000000940563409130278003002941",
			shouldProgress: false,
			expecting:      Step{},
		},
		{
			name:           "HiddenTriple",
			technique:      HiddenTriple,
			board:          "385601000109500000020030510000005000030010060000400000017050080003100900000007132",
			shouldProgress: true,
			expecting: Step{
				Technique:         "HiddenTriple",
				ReasonCells:       NoPeers.Including([]board.Coordinates{{Row: 1, Col: 6}, {Row: 1, Col: 8}, {Row: 2, Col: 8}}...),
				AffectedCells:     NoPeers.Including([]board.Coordinates{{Row: 1, Col: 6}, {Row: 1, Col: 8}, {Row: 2, Col: 8}}...),
				RemovedCandidates: board.MustCandidateSet(2, 4, 7, 9),
			},
		},
		{
			name:           "HiddenTriple (No Progress)",
			technique:      HiddenTriple,
			board:          "975421386148563700632879154006200430004300000390004000000940563409130278003002941",
			shouldProgress: false,
			expecting:      Step{},
		},
		{
			name:           "HiddenQuad",
			technique:      HiddenQuad,
			board:          "000000001000010002000634905000309056003860049679145823000251638126783594835496217",
			shouldProgress: true,
			expecting: Step{
				Technique:         "HiddenQuad",
				ReasonCells:       NoPeers.Including([]board.Coordinates{{Row: 0, Col: 0}, {Row: 0, Col: 1}, {Row: 1, Col: 0}, {Row: 1, Col: 1}}...),
				AffectedCells:     NoPeers.Including([]board.Coordinates{{Row: 0, Col: 0}, {Row: 0, Col: 1}, {Row: 1, Col: 0}, {Row: 1, Col: 1}}...),
				RemovedCandidates: board.MustCandidateSet(2, 4, 7, 8),
			},
		},
		{
			name:           "HiddenQuad (No Progress)",
			technique:      HiddenQuad,
			board:          "975421386148563700632879154006200430004300000390004000000940563409130278003002941",
			shouldProgress: false,
			expecting:      Step{},
		},
		{
			name:           "LockedCandidates",
			technique:      LockedCandidates,
			board:          "984000000002500040001904002006097230003602000209035610195768423427351896638009751",
			shouldProgress: true,
			expecting: Step{
				Technique:         "LockedCandidates (Row)",
				ReasonCells:       NoPeers.Including([]board.Coordinates{{Row: 0, Col: 6}, {Row: 0, Col: 8}}...),
				AffectedCells:     NoPeers.Including([]board.Coordinates{{Row: 2, Col: 6}}...),
				RemovedCandidates: board.MustCandidateSet(5),
			},
		},
		{
			name:           "LockedCandidates (No Progress)",
			technique:      LockedCandidates,
			board:          "010000496060941807490060100129534678840796010706128940084610709900480061601079084",
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
			for affected := range step.AffectedCells.Each() {
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
