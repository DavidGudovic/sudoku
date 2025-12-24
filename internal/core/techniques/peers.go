package techniques

import (
	"iter"
	"math/bits"

	"github.com/DavidGudovic/sudoku/internal/core/board"
)

//go:generate stringer -type=Scope

var (
	NoPeers          = PeerSet{}
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

// Peers is the entry point for building PeerSet's based on board.Coordinates and Scope's.
var Peers peerQuery

// WithCoordinates is a helper type for building PeerSet's based on board.Coordinates.
type WithCoordinates struct {
	coords []board.Coordinates
}

// Of starts building a PeerSet for the specified board.Coordinates.
//
// To finalize the query and build the PeerSet, call Across(...Scope) or AcrossSharedScopes().
func (peerQuery) Of(coords ...board.Coordinates) WithCoordinates {
	return WithCoordinates{coords: coords}
}

// All finalizes the query with a PeerSet containing all board.Coordinates on the board.
func (peerQuery) All() PeerSet {
	return PeerSet{
		0b111111111,
		0b111111111,
		0b111111111,
		0b111111111,
		0b111111111,
		0b111111111,
		0b111111111,
		0b111111111,
		0b111111111,
	}
}

// Empty finalizes the query with an empty PeerSet.
func (peerQuery) Empty() PeerSet {
	return NoPeers
}

// InScope finalizes the query with a PeerSet containing all board.Coordinates in the specified Scope and index.
func (peerQuery) InScope(scope Scope, index int) PeerSet {
	ps := NoPeers

	switch scope {
	case Row:
		ps[index] = FullRow
	case Column:
		for r := 0; r < board.Size; r++ {
			ps[r] |= 1 << index
		}
	case Box:
		for i := 0; i < board.BoxSize*board.BoxSize; i++ {
			ps = ps.With(board.MustCoordsFromBoxIndex(index, i))
		}
	}

	return ps
}

// Across builds a PeerSet containing all peers of the specified board.Coordinates across the given scopes.
func (w WithCoordinates) Across(scopes ...Scope) PeerSet {
	var ps PeerSet

	for _, c := range w.coords {
		for _, s := range scopes {
			switch s {
			case Row:
				ps[c.Row] = FullRow
			case Column:
				colMask := uint16(1 << c.Col)
				for r := 0; r < board.Size; r++ {
					ps[r] |= colMask
				}
			case Box:
				for i := 0; i < board.Size; i++ {
					ps = ps.With(board.MustCoordsFromBoxIndex(c.BoxIndex(), i))
				}
			}
		}
	}

	return ps.Excluding(w.coords...)
}

// AcrossSharedScopes builds a PeerSet containing all peers of the specified board.Coordinates across shared Scope.
// Shared Scope is defined as a scope where all coordinates share the same row, column, or box.
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
		scope     Scope
		predicate func(board.Coordinates) bool
	}{
		{Row, first.SharesRowWith},
		{Column, first.SharesColumnWith},
		{Box, first.SharesBoxWith},
	}

	var scopes []Scope
	for _, chk := range checks {
		if allCoords(coords[1:], chk.predicate) {
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

// Union returns a new PeerSet that is the union of this set and another.
func (ps PeerSet) Union(other PeerSet) PeerSet {
	result := NoPeers

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
		coords := board.MustCoordsFromBoxIndex(boxIndex, i)
		if ps.Contains(coords) {
			return true
		}
	}
	return false
}

// Count returns the amount of board.Coordinates in the set.
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
	result := NoPeers

	for r := 0; r < board.Size; r++ {
		result[r] = ps[r] &^ other[r]
	}

	return result
}

// Intersection returns a new PeerSet that is the intersection of this set and another.
func (ps PeerSet) Intersection(other PeerSet) PeerSet {
	result := NoPeers

	for r := 0; r < board.Size; r++ {
		result[r] = ps[r] & other[r]
	}

	return result
}

// IsEmpty checks if the PeerSet is empty.
func (ps PeerSet) IsEmpty() bool {
	return ps == NoPeers
}

// Candidates return the union of candidates from all board.Cell's in this PeerSet on a given board.Board.
func (ps PeerSet) Candidates(p board.Board) board.CandidateSet {
	seen := board.NoCandidates

	for c := range ps.Each() {
		seen.Merge(p.CellAt(c).Candidates())
	}

	return seen
}

// ContainingCandidates filters the PeerSet to only include cells that contain any of the specified candidates on the given board.Board.
func (ps PeerSet) ContainingCandidates(p board.Board, candidates board.CandidateSet) PeerSet {
	result := NoPeers

	for c := range ps.Each() {
		set := p.CellAt(c).Candidates()

		if set.Intersection(candidates) != board.NoCandidates {
			result[c.Row] |= 1 << c.Col
		}
	}

	return result
}

// NotContainingCandidates filters the PeerSet to only include cells that do not contain the specified candidates on the given board.Board.
func (ps PeerSet) NotContainingCandidates(p board.Board, candidates board.CandidateSet) PeerSet {
	result := NoPeers

	for c := range ps.Each() {
		set := p.CellAt(c).Candidates()

		if set.Intersection(candidates) == board.NoCandidates {
			result[c.Row] |= 1 << c.Col
		}
	}

	return result
}

// ContainingExactCandidates filters the PeerSet to only include cells that contain exactly the specified candidates on the given board.Board.
func (ps PeerSet) ContainingExactCandidates(p board.Board, candidates board.CandidateSet) PeerSet {
	result := NoPeers

	for c := range ps.Each() {
		target := p.CellAt(c).Candidates()

		if target == candidates {
			result = result.With(c)
		}
	}

	return result
}

// EmptyCells filters the PeerSet to only include cells that are empty on the given board.Board.
func (ps PeerSet) EmptyCells(p board.Board) PeerSet {
	result := NoPeers

	for c := range ps.Each() {
		if p.CellAt(c).IsEmpty() {
			result = result.With(c)
		}
	}

	return result
}

// ContainingValues filters the PeerSet to only include cells that contain any of the specified values on the given board.Board.
func (ps PeerSet) ContainingValues(p board.Board, values ...int) PeerSet {
	result := NoPeers

	valueSet := board.MustCandidateSet(values...)
	for c := range ps.Each() {
		if valueSet.Contains(p.CellAt(c).Value()) {
			result = result.With(c)
		}
	}

	return result
}

// Each executes a function for every coordinate present in the PeerSet.
// It uses bit-scanning to skip empty cells, making it much faster than 9x9 loops.
func (ps PeerSet) Each() iter.Seq[board.Coordinates] {
	return func(yield func(board.Coordinates) bool) {
		for r := 0; r < board.Size; r++ {
			mask := ps[r]
			for mask != 0 {
				c := bits.TrailingZeros16(mask)
				if !yield(board.MustCoordinates(r, c)) {
					return
				}
				mask &^= 1 << c
			}
		}
	}
}

// Slice converts to a slice of board.Coordinates.
// Used for when the caller needs a slice, e.g. for Step.AffectedCells.
// Should not be used for iteration, instead use Each()
func (ps PeerSet) Slice() []board.Coordinates {
	var coords []board.Coordinates

	for row := 0; row < board.Size; row++ {
		for col := 0; col < board.Size; col++ {
			c := board.MustCoordinates(row, col)
			if ps.Contains(c) {
				coords = append(coords, c)
			}
		}
	}

	return coords
}
