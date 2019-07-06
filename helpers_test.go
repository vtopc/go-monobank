package monobank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToBanknote(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name       string
		i          int64
		minorUnits int
		want       string
	}{
		{
			name:       "positive",
			i:          123,
			minorUnits: 2,
			want:       "1.23",
		},
		{
			name:       "less_then_one_banknote-indent",
			i:          34,
			minorUnits: 2,
			want:       "0.34",
		},
		{
			name:       "indent-2",
			i:          4,
			minorUnits: 2,
			want:       "0.04",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToBanknote(tt.i, tt.minorUnits)
			assert.Equal(t, tt.want, got)
		})
	}
}
