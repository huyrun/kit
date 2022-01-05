package env

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func TestGetWithDefault(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		key string
		value any
	}{
		{"string", "test_string", "value"},
		{"int", "test_int", rand.Int()},
		{"time", "test_time", time.Now()},
	}
	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := GetWithDefault(tc.key, tc.value)
			require.Equal(t, tc.value, got)
		})
	}
}
