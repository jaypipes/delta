package genutil_test

import (
	"testing"

	"github.com/jaypipes/delta/internal/genutil"
	"github.com/stretchr/testify/require"
)

func TestUniquePairs(t *testing.T) {
	require := require.New(t)
	subject := []any{
		int(0),
		uint(0),
		float32(0),
	}
	exp := [][2]any{
		{int(0), uint(0)},
		{int(0), float32(0)},
		{uint(0), int(0)},
		{uint(0), float32(0)},
		{float32(0), int(0)},
		{float32(0), uint(0)},
	}

	got := genutil.UniquePairs(subject)
	require.Equal(exp, got)
}
