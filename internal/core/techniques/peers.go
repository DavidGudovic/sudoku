package techniques

import (
	"math/bits"

	"github.com/DavidGudovic/sudoku/internal/core/board"
)

var (
	NoCells          = PeerSet{}
	FullRow   uint16 = 0b111111111
	AllScopes        = []Scope{Row, Column, Box}
)

const (
	Row Scope = iota
	Column
	Box
)

// PeerSet is an array of bit masks representing board.Coordinates.
// Each element represents one row, with 9 bits for the 9 columns.
//
// For example, if PeerSet[0] = 0b000000101, it means that in row 0, columns 0 and 2 are included in the set.
// This is used for efficient representation and manipulation of sets of board.Coordinates for techniques.
type PeerSet [board.Size]uint16

// Scope represents a unit type in Sudoku: Row, Column, or Box. (also known as a "house", or "unit")
type Scope int

type peerQuery struct{}

// Peers is the entry point for building PeerSet's based on board.Coordinates and scopes.
var Peers peerQuery

// Of starts building a PeerSet for the specified board.Coordinates.
// To finalize the query and build the PeerSet, call Across(...Scope) or AcrossSharedScopes().
func (peerQuery) Of(coords ...board.Coordinates) WithCoordinates {
	return WithCoordinates{coords: coords}
}

type WithCoordinates struct {
	coords []board.Coordinates
}

// Across builds a PeerSet containing all peers of the specified board.Coordinates across the given scopes.
func (w WithCoordinates) Across(scopes ...Scope) PeerSet {
	var ps PeerSet

	for _, c := range w.coords {
		for _, s := range scopes {
			switch s {
			case Row:
				ps = ps.WithRow(c.Row)
			case Column:
				ps = ps.WithCol(c.Col)
			case Box:
				ps = ps.WithBox(c.BoxIndex())
			}
		}
	}

	return ps.Excluding(w.coords...)
}

// AcrossSharedScopes builds a PeerSet containing all peers of the specified board.Coordinates across shared scopes.
func (w WithCoordinates) AcrossSharedScopes() PeerSet {
	return w.Across(SharedScopesOf(w.coords)...)
}

// SharedScopesOf returns a slice of Scope representing the shared scopes between all board.Coordinates.
func SharedScopesOf(coords []board.Coordinates) []Scope {
	if len(coords) == 0 {
		return nil
	}

	first := coords[0]
	checks := []struct {
		scope Scope
		check func(board.Coordinates) bool
	}{
		{Row, first.SharesRowWith},
		{Column, first.SharesColumnWith},
		{Box, first.SharesBoxWith},
	}

	var scopes []Scope
	for _, chk := range checks {
		if allCoords(coords[1:], chk.check) {
			scopes = append(scopes, chk.scope)
		}
	}

	return scopes
}

// allCoords checks if all coordinates satisfy the given predicate.
func allCoords(coords []board.Coordinates, predicate func(board.Coordinates) bool) bool {
	for _, c := range coords {
		if !predicate(c) {
			return false
		}
	}
	return true
}

// Contains checks if the set contains the specified board.Coordinates.
func (ps PeerSet) Contains(c board.Coordinates) bool {
	return ps[c.Row]&(1<<c.Col) != 0
}

// With returns a new PeerSet with the specified board.Coordinates added.
func (ps PeerSet) With(c board.Coordinates) PeerSet {
	result := ps
	result[c.Row] |= 1 << c.Col
	return result
}

// Without returns a new PeerSet with the specified board.Coordinates removed.
func (ps PeerSet) Without(c board.Coordinates) PeerSet {
	result := ps
	result[c.Row] &^= 1 << c.Col
	return result
}

// WithRow returns a new PeerSet with all cells in the specified row added.
func (ps PeerSet) WithRow(row int) PeerSet {
	result := ps
	result[row] = FullRow
	return result
}

// WithCol returns a new PeerSet with all cells in the specified column added.
func (ps PeerSet) WithCol(col int) PeerSet {
	result := ps

	for r := 0; r < board.Size; r++ {
		result[r] |= 1 << col
	}

	return result
}

// WithBox returns a new PeerSet with all cells in the specified box added.
func (ps PeerSet) WithBox(boxIndex int) PeerSet {
	result := ps

	for i := 0; i < board.BoxSize*board.BoxSize; i++ {
		coords, _ := board.CoordsFromBoxIndex(boxIndex, i)
		result = result.With(coords)
	}

	return result
}

// Union returns a new PeerSet that is the union of this set and another.
func (ps PeerSet) Union(other PeerSet) PeerSet {
	result := NoCells

	for r := 0; r < board.Size; r++ {
		result[r] = ps[r] | other[r]
	}

	return result
}

// HasPeersInRow checks if the PeerSet contains any peers in the specified row.
func (ps PeerSet) HasPeersInRow(row int) bool {
	return ps[row] != 0
}

// HasPeersInCol checks if the PeerSet contains any peers in the specified column.
func (ps PeerSet) HasPeersInCol(col int) bool {
	for r := 0; r < board.Size; r++ {
		if ps[r]&(1<<col) != 0 {
			return true
		}
	}
	return false
}

// HasPeersInBox checks if the PeerSet contains any peers in the specified box.
func (ps PeerSet) HasPeersInBox(boxIndex int) bool {
	for i := 0; i < board.BoxSize*board.BoxSize; i++ {
		coords, _ := board.CoordsFromBoxIndex(boxIndex, i)
		if ps.Contains(coords) {
			return true
		}
	}
	return false
}

// PeersInRow returns a new PeerSet containing only the peers in the specified row.
func (ps PeerSet) PeersInRow(row int) PeerSet {
	result := NoCells
	result[row] = ps[row]
	return result
}

// PeersInCol returns a new PeerSet containing only the peers in the specified column.
func (ps PeerSet) PeersInCol(col int) PeerSet {
	result := NoCells
	for r := 0; r < board.Size; r++ {
		result[r] = ps[r] & (1 << col)
	}
	return result
}

// PeersInBox returns a new PeerSet containing only the peers in the specified box.
func (ps PeerSet) PeersInBox(boxIndex int) PeerSet {
	result := NoCells
	for i := 0; i < board.BoxSize*board.BoxSize; i++ {
		coords, _ := board.CoordsFromBoxIndex(boxIndex, i)
		if ps.Contains(coords) {
			result = result.With(coords)
		}
	}
	return result
}

// Count returns the number of board.Coordinates in the set.
func (ps PeerSet) Count() int {
	count := 0
	for _, row := range ps {
		count += bits.OnesCount16(row)
	}
	return count
}

// Including returns a new PeerSet with the specified board.Coordinates added.
func (ps PeerSet) Including(coords ...board.Coordinates) PeerSet {
	result := ps
	for _, c := range coords {
		result = result.With(c)
	}
	return result
}

// Excluding returns a new PeerSet with the specified board.Coordinates removed.
func (ps PeerSet) Excluding(coords ...board.Coordinates) PeerSet {
	result := ps
	for _, c := range coords {
		result = result.Without(c)
	}
	return result
}

// Except returns a new PeerSet excluding the given PeerSet
func (ps PeerSet) Except(other PeerSet) PeerSet {
	result := NoCells

	for r := 0; r < board.Size; r++ {
		result[r] = ps[r] &^ other[r]
	}

	return result
}

// Intersection returns a new PeerSet that is the intersection of this set and another.
func (ps PeerSet) Intersection(other PeerSet) PeerSet {
	result := NoCells

	for r := 0; r < board.Size; r++ {
		result[r] = ps[r] & other[r]
	}

	return result
}

// IsEmpty checks if the PeerSet is empty.
func (ps PeerSet) IsEmpty() bool {
	return ps == NoCells
}

// Candidates return the union of candidates from all board.Cell's in this PeerSet on a given board.Board.
func (ps PeerSet) Candidates(p board.Board) board.CandidateSet {
	seen := board.NoCandidates

	for row := 0; row < board.Size; row++ {
		for col := 0; col < board.Size; col++ {
			if ps[row]&(1<<col) != 0 {
				coords, _ := board.NewCoordinates(row, col)
				seen.Merge(p.CellAt(coords).Candidates())
			}
		}
	}

	return seen
}

// ContainingCandidates filters the PeerSet to only include cells that contain any of the specified candidates on the given board.Board.
func (ps PeerSet) ContainingCandidates(p board.Board, candidates board.CandidateSet) PeerSet {
	result := NoCells

	coords := ps.Slice()

	for _, c := range coords {
		cellCandidates := p.CellAt(c).Candidates()

		if cellCandidates.Intersection(candidates) != board.NoCandidates {
			result = result.With(c)
		}
	}

	return result
}

// NotContainingCandidates filters the PeerSet to only include cells that do not contain the specified candidates on the given board.Board.
func (ps PeerSet) NotContainingCandidates(p board.Board, candidates board.CandidateSet) PeerSet {
	result := NoCells

	coords := ps.Slice()
	for _, c := range coords {
		cellCandidates := p.CellAt(c).Candidates()

		if cellCandidates.Intersection(candidates) == board.NoCandidates {
			result = result.With(c)
		}
	}

	return result
}

// ContainingExactCandidates filters the PeerSet to only include cells that contain exactly the specified candidates on the given board.Board.
func (ps PeerSet) ContainingExactCandidates(p board.Board, candidates board.CandidateSet) PeerSet {
	result := NoCells

	coords := ps.Slice()
	for _, c := range coords {
		target := p.CellAt(c).Candidates()

		if target == candidates {
			result = result.With(c)
		}
	}

	return result
}

// Slice converts to a slice of board.Coordinates.
func (ps PeerSet) Slice() []board.Coordinates {
	var coords []board.Coordinates

	for row := 0; row < board.Size; row++ {
		for col := 0; col < board.Size; col++ {
			c, _ := board.NewCoordinates(row, col)
			if ps.Contains(c) {
				coords = append(coords, c)
			}
		}
	}

	return coords
}
