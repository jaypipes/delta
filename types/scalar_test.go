package types_test

import (
	"fmt"
	"slices"
	"sort"
	"testing"

	testtypes "github.com/jaypipes/delta/internal/testing/types"
	"github.com/jaypipes/delta/types"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestIsScalar(t *testing.T) {
	notScalars := []any{
		testtypes.NilInterface,
		testtypes.NotNilInterface,
		struct{ a string }{a: "1"},
		func() {},
		nil,
		map[string]string{"a": "a"},
		[]string{"a"},
	}
	t.Run("not scalars", func(t *testing.T) {
		require := require.New(t)
		for _, subject := range notScalars {
			require.False(
				types.IsScalar(subject),
				fmt.Sprintf(
					"expected %v (%T) not to be a scalar",
					subject, subject,
				),
			)
		}
	})

	scalars := slices.Concat(
		// All primitives are scalars
		testtypes.PrimitiveZeroValues,
		// Type definitions over built-in primitives are scalars.
		testtypes.PrimitiveTypeDefZeroValues,
	)

	t.Run("scalars", func(t *testing.T) {
		require := require.New(t)
		for _, subject := range scalars {
			require.True(
				types.IsScalar(subject),
				fmt.Sprintf(
					"expected %v (%T) to be a scalar",
					subject, subject,
				),
			)
		}
	})
}

func ExampleIsScalar() {

	m := map[string]any{
		"nil":               nil,
		"nil iface":         testtypes.NilInterface,
		"struct":            testtypes.MyStructType{},
		"func":              func() {},
		"string":            "1",
		"[]string":          []string{"a", "b"},
		"map[string]string": map[string]string{"a": "b"},
		"[]byte":            []byte("1"),
		"rune":              rune('1'),
		"int32":             int32('1'),
		"byte":              byte('1'),
		"uint8":             uint8('1'),
		"int":               1,
		"uint":              uint(1),
		"float32":           float32(1),
		"float64":           float64(1),
		"bool":              true,
	}
	keys := lo.Keys(m)
	sort.Strings(keys)
	for _, k := range keys {
		v := m[k]
		fmt.Printf("Is %s scalar? %t\n", k, types.IsScalar(v))
	}

	// Output:
	// Is []byte scalar? true
	// Is []string scalar? false
	// Is bool scalar? true
	// Is byte scalar? true
	// Is float32 scalar? true
	// Is float64 scalar? true
	// Is func scalar? false
	// Is int scalar? true
	// Is int32 scalar? true
	// Is map[string]string scalar? false
	// Is nil scalar? false
	// Is nil iface scalar? false
	// Is rune scalar? true
	// Is string scalar? true
	// Is struct scalar? false
	// Is uint scalar? true
	// Is uint8 scalar? true
}
