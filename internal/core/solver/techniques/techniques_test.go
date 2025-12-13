package techniques

import "testing"

func TestTechniques(t *testing.T) {
	tests := []struct {
		name      string
		technique Technique
	}{
		{
			name:      "NakedSingle",
			technique: NakedSingle{},
		},
	}

	_ = tests
}
