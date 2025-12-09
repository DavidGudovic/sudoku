package board

import "strings"

func CoordsFromIndex(index int) (int, int, error) {
	if index < 0 || index >= CellCount {
		return 0, 0, ErrIndexOutOfBounds
	}

	return index / Size, index % Size, nil
}

func filterCandidates(s string) (string, error) {
	var sb strings.Builder
	isCandidate := false

	for _, ch := range s {
		if isCandidate {
			isCandidate = false
			continue
		}

		if ch == CandidatePrefixRune {
			isCandidate = true
			continue
		}

		if ch < EmptyCell+'0' || ch > MaxValue+'0' {
			return "", ErrInvalidRuneInStringRep
		}

		sb.WriteRune(ch)
	}

	if isCandidate {
		return "", ErrInvalidStringRep
	}

	return sb.String(), nil
}
