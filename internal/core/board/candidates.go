package board

import (
	"iter"
	"math/bits"
	"strconv"
	"strings"
)

var (
	AllCandidates CandidateSet = 0b1111111110
	NoCandidates  CandidateSet = 0b0000000000
)

// CandidateSet is a bit mask representing the possible candidates for a Cell.
// The least significant bit (bit 0) is unused, bits 1-9 represent candidates 1-9.
//   - For example, a CandidateSet of 0b0000100110 represents candidates 1, 2, and 5.
type CandidateSet uint16

// NewCandidateSet creates a new CandidateSet from the provided values.
func NewCandidateSet(values ...int) (CandidateSet, error) {
	var cs CandidateSet

	for _, val := range values {
		err := cs.Add(val)
		if err != nil {
			return NoCandidates, err
		}
	}

	return cs, nil
}

// MustCandidateSet creates a new CandidateSet from the provided values, assuming the caller guarantees valid input.
// If the input is invalid, there's either a serious bug in the caller or the world view is wrong, therefore it panics.
func MustCandidateSet(values ...int) CandidateSet {
	cs, err := NewCandidateSet(values...)

	if err != nil {
		panic("Impossible: " + err.Error())
	}

	return cs
}

// Contains checks if the CandidateSet contains the specified candidate value.
func (cs *CandidateSet) Contains(value int) bool {
	return *cs&(1<<value) != 0
}

// Add adds (sets the bit to 1) the specified candidate value to the CandidateSet.
func (cs *CandidateSet) Add(value int) error {
	if value < MinValue || value > MaxValue {
		return ErrInvalidCellValue
	}

	*cs |= 1 << value

	return nil
}

// Remove removes (sets the bit to 0) the specified candidate value from the CandidateSet.
func (cs *CandidateSet) Remove(value int) {
	*cs &^= 1 << value
}

// Exclude removes all candidates present in the other CandidateSet from the current CandidateSet.
func (cs *CandidateSet) Exclude(other CandidateSet) {
	*cs &^= other
}

// Merge adds all candidates present in the other CandidateSet to the current CandidateSet.
func (cs *CandidateSet) Merge(other CandidateSet) {
	*cs |= other
}

// Intersection returns a new CandidateSet that is the intersection of the current CandidateSet and another CandidateSet.
func (cs *CandidateSet) Intersection(other CandidateSet) CandidateSet {
	return *cs & other
}

// Count returns the number of candidates in a CandidateSet
func (cs *CandidateSet) Count() int {
	return bits.OnesCount16(uint16(*cs))
}

// Slice converts the CandidateSet to a slice of integers representing the candidate values.
func (cs *CandidateSet) Slice() []int {
	var values []int

	for val := range cs.Each() {
		values = append(values, val)
	}

	return values
}

// Each executes a function for every candidate present in the CandidateSet.
// It uses bit-scanning to skip empty candidates.
func (cs *CandidateSet) Each() iter.Seq[int] {
	return func(yield func(int) bool) {
		mask := *cs
		for mask != 0 {
			c := bits.TrailingZeros16(uint16(mask))
			if !yield(c) {
				return
			}
			mask &^= 1 << c
		}
	}
}

// First returns the first candidate value in the CandidateSet 1 - 9, or 0 if empty.
func (cs *CandidateSet) First() int {
	return bits.TrailingZeros16(uint16(*cs))
}

// String converts the CandidateSet to a comma-separated list of integer values of the candidates.
func (cs *CandidateSet) String() string {
	sb := strings.Builder{}

	for val := range cs.Each() {
		if sb.Len() > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.Itoa(val))
	}

	return sb.String()
}

// CandidateSubsets generates all possible CandidateSet combinations of size n from the given CandidateSet.
// Mathematically, this is equivalent to computing "n choose k" (nCk) combinations.
func CandidateSubsets(cs CandidateSet, size int) []CandidateSet {
	if size == 0 {
		return []CandidateSet{NoCandidates}
	}

	if cs.Count() < size {
		return nil
	}

	var first CandidateSet
	for c := range cs.Each() {
		_ = first.Add(c)
		break
	}

	cs.Remove(first.First())

	rest := cs

	var results []CandidateSet
	for _, combo := range CandidateSubsets(rest, size-1) {
		combo.Merge(first)
		results = append(results, combo)
	}

	results = append(results, CandidateSubsets(rest, size)...)

	return results
}
