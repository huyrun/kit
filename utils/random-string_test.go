package utils

import (
	"github.com/huypher/kit"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandomString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		n int
		alphabet string
	}{
		{"success", 10, kit.EnglishLetters},
	}
	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := RandomString(tc.n, tc.alphabet)
			require.NotNil(t, got)
		})
	}
}
