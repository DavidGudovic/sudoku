package techniques

import (
	"github.com/DavidGudovic/sudoku/internal/core/board"
)

// Technique represents a Sudoku solving technique that can be applied to a puzzle
type Technique interface {
	Apply(puzzle *board.Board) (Step, error)
}

// Func is a function type that implements the Technique interface
// This allows simple functions to be used as techniques
type Func func(*board.Board) (Step, error)

func (tf Func) Apply(puzzle *board.Board) (Step, error) {
	return tf(puzzle)
}

// Step represents a single solving step taken by a technique
type Step struct {
	Description       string
	Technique         string
	AffectedCells     []board.Coordinates
	ReasonCells       []board.Coordinates
	RemovedCandidates board.CandidateSet
	PlacedValue       *int
}

// StepStack represents a stack of Steps, typically used to track the sequence of solving steps
type StepStack []Step

func (s Step) MadeProgress() bool {
	return s.PlacedValue != nil || s.RemovedCandidates != board.NoCandidates
}

// PeerCoordinates represents a slice of board coordinates that are peers to a given cell
type PeerCoordinates []board.Coordinates

// Excluding removes the given coordinates from the PeerCoordinates slice
func (p PeerCoordinates) Excluding(coords board.Coordinates) PeerCoordinates {
	var result PeerCoordinates
	for _, c := range p {
		if c != coords {
			result = append(result, c)
		}
	}
	return result
}

// Including adds the given coordinates to the PeerCoordinates slice
func (p PeerCoordinates) Including(coords board.Coordinates) PeerCoordinates {
	return append(p, coords)
}

// rowPeers returns all peer coordinates in the given row
func rowPeers(r int) PeerCoordinates {
	var pc PeerCoordinates

	for i := 0; i < board.Size; i++ {
		coords, _ := board.NewCoordinates(r, i)
		pc = append(pc, coords)
	}

	return pc
}

// colPeers returns all peer coordinates in the given column index (0-8)
func colPeers(c int) PeerCoordinates {
	var pc PeerCoordinates

	for i := 0; i < board.Size; i++ {
		coords, _ := board.NewCoordinates(i, c)
		pc = append(pc, coords)
	}

	return pc
}

// boxPeers returns all peer coordinates in the given box index (0-8)
func boxPeers(b int) PeerCoordinates {
	var pc PeerCoordinates

	for i := 0; i < board.BoxSize*board.BoxSize; i++ {
		coords, _ := board.CoordsFromBoxIndex(b, i)
		pc = append(pc, coords)
	}

	return pc
}

// allPeers returns all peer coordinates for the given cell coordinates
func allPeers(c board.Coordinates) PeerCoordinates {
	return append(append(rowPeers(c.Row), colPeers(c.Col)...), boxPeers(c.BoxIndex())...)
}

// rowCandidates returns a CandidateSet representing all candidates present in the given row at least once
func rowCandidates(p board.Board, row int) board.CandidateSet {
	seen := board.NoCandidates

	for col := 0; col < board.Size; col++ {
		cell := p.CellAt(board.Coordinates{Row: row, Col: col})
		seen.Merge(cell.Candidates())
	}

	return seen
}

// columnCandidates returns a CandidateSet representing all candidates present in the given column at least once
func columnCandidates(p board.Board, col int) board.CandidateSet {
	seen := board.NoCandidates

	for row := 0; row < board.Size; row++ {
		cell := p.CellAt(board.Coordinates{Row: row, Col: col})
		seen.Merge(cell.Candidates())
	}

	return seen
}

// boxCandidates returns a CandidateSet representing all candidates present in the given box at least once
func boxCandidates(p board.Board, coords board.Coordinates) board.CandidateSet {
	seen := board.NoCandidates

	boxIndex := coords.BoxIndex()

	for i := 0; i < board.BoxSize*board.BoxSize; i++ {
		bc, _ := board.CoordsFromBoxIndex(boxIndex, i)
		cell := p.CellAt(bc)
		seen.Merge(cell.Candidates())
	}

	return seen
}
