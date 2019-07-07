package monobank

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatementUnmarshal(t *testing.T) {
	ts := int64(1554466347)

	tests := []struct {
		name    string
		v       string
		want    Statement
		wantErr error
	}{
		{
			v: fmt.Sprintf(`{"time": %d}`, ts),
			want: Statement{
				Time: Time{time.Unix(ts, 0)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Statement
			err := json.Unmarshal([]byte(tt.v), &got)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
