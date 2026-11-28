package types_test

import (
	"fmt"
	"slices"
	"sort"
	"testing"

	testtypes "github.com/jaypipes/delta/internal/testing/types"
	"github.com/jaypipes/delta/types"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
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
		testtypes.MyStructType{A: 1},
		[]testtypes.MyStructType{testtypes.MyStructType{A: 1}},
		[]testtypes.MyIntType{1},
		testtypes.MySliceIntType{1},
	}
	t.Run("not scalars", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range notScalars {
			assert.False(
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
		assert := assert.New(t)
		for _, subject := range scalars {
			assert.True(
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

	type myInterfaceType interface {
		F()
	}

	type myStructType struct {
		myInterfaceType
		A int
	}

	var nilInterface myInterfaceType = (*myStructType)(nil)

	type myStringType string
	type myByteArrayType []byte
	type myInt32Type int
	type myFloat32Type float32
	type myFloat64Type float64

	type mySliceStringType []string
	type mySliceIntType []int
	type mySliceFloat32Type []float32
	type mySliceFloat64Type []float64

	m := map[string]any{
		"nil":               nil,
		"nil iface":         nilInterface,
		"struct":            myStructType{},
		"func":              func() {},
		"string":            "foo",
		"string type def":   myStringType("foo"),
		"[]string":          []string{"a", "b"},
		"[]string type def": mySliceStringType{"a", "b"},
		"map[string]string": map[string]string{"a": "b"},
		"[]byte":            []byte("1"),
		"[]byte type def":   myByteArrayType("1"),
		"rune":              rune('1'),
		"int32":             int32('1'),
		"int32 type def":    myInt32Type('1'),
		"byte":              byte('1'),
		"uint8":             uint8('1'),
		"int":               1,
		"uint":              uint(1),
		"float32":           float32(1),
		"float32 type def":  myFloat32Type(1),
		"float64":           float64(1),
		"float64 type def":  myFloat64Type(1),
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
	// Is []byte type def scalar? true
	// Is []string scalar? false
	// Is []string type def scalar? false
	// Is bool scalar? true
	// Is byte scalar? true
	// Is float32 scalar? true
	// Is float32 type def scalar? true
	// Is float64 scalar? true
	// Is float64 type def scalar? true
	// Is func scalar? false
	// Is int scalar? true
	// Is int32 scalar? true
	// Is int32 type def scalar? true
	// Is map[string]string scalar? false
	// Is nil scalar? false
	// Is nil iface scalar? false
	// Is rune scalar? true
	// Is string scalar? true
	// Is string type def scalar? true
	// Is struct scalar? false
	// Is uint scalar? true
	// Is uint8 scalar? true
}
