package delta_test

import (
	"testing"

	"github.com/jaypipes/delta"
	"github.com/stretchr/testify/require"
)

type iface interface {
	f()
}
type s struct {
	a int
}

func (s *s) f() {}

var (
	sNil    iface = (*s)(nil)
	sNotNil iface = &s{1}
)

func TestScalar_IsScalar(t *testing.T) {
	cases := []struct {
		name    string
		subject any
		exp     bool
	}{
		{
			"interface with nil underlying value is not scalar",
			sNil,
			false,
		},
		{
			"interface with non-nil underlying value is not scalar",
			sNotNil,
			false,
		},
		{
			"slice of scalar is not scalar",
			[]string{"1"},
			false,
		},
		{
			"struct with scalar is not scalar",
			struct{ a string }{a: "1"},
			false,
		},
		{
			"func is not scalar",
			func() {},
			false,
		},
		{
			"nil is scalar",
			nil,
			true,
		},
		{
			"int is scalar",
			int(0),
			true,
		},
		{
			"int8 is scalar",
			int8(0),
			true,
		},
		{
			"int16 is scalar",
			int16(0),
			true,
		},
		{
			"int32 is scalar",
			int32(0),
			true,
		},
		{
			"int64 is scalar",
			int64(0),
			true,
		},
		{
			"uint is scalar",
			uint(0),
			true,
		},
		{
			"uint8 is scalar",
			uint8(0),
			true,
		},
		{
			"uint16 is scalar",
			uint16(0),
			true,
		},
		{
			"uint32 is scalar",
			uint32(0),
			true,
		},
		{
			"uint64 is scalar",
			uint64(0),
			true,
		},
		{
			"string is scalar",
			"foor",
			true,
		},
		{
			"[]byte is scalar",
			[]byte("foor"),
			true,
		},
		{
			"rune is scalar",
			'c',
			true,
		},
		{
			"float32 is scalar",
			float32(0),
			true,
		},
		{
			"float64 is scalar",
			float64(0),
			true,
		},
		{
			"bool is scalar",
			false,
			true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(c.exp, delta.IsScalar(c.subject))
		})
	}
}

func TestScalar_ScalarEquivalent(t *testing.T) {
	cases := []struct {
		name string
		a    any
		b    any
		exp  bool
	}{
		{
			"two nils are equivalent",
			nil,
			nil,
			true,
		},
		{
			"true and \"true\" are equivalent",
			true,
			"true",
			true,
		},
		{
			"true and \"t\" are equivalent",
			true,
			"t",
			true,
		},
		{
			"true and \"1\" are equivalent",
			true,
			"1",
			true,
		},
		{
			"true and []byte(\"true\") are equivalent",
			true,
			[]byte("true"),
			true,
		},
		{
			"true and 1 are equivalent",
			true,
			1,
			true,
		},
		{
			"true and uint(1) are equivalent",
			true,
			uint(1),
			true,
		},
		{
			"true and int16(1) are equivalent",
			true,
			int16(1),
			true,
		},
		{
			"true and float32(1) are not equivalent",
			true,
			float32(1),
			false,
		},
		{
			"true and float64(1) are not equivalent",
			true,
			float64(1),
			false,
		},
		{
			"false and \"false\" are equivalent",
			false,
			"false",
			true,
		},
		{
			"false and \"f\" are equivalent",
			false,
			"f",
			true,
		},
		{
			"false and \"0\" are equivalent",
			false,
			"0",
			true,
		},
		{
			"false and []byte(\"false\") are equivalent",
			false,
			[]byte("false"),
			true,
		},
		{
			"false and 0 are equivalent",
			false,
			0,
			true,
		},
		{
			"false and uint(0) are equivalent",
			false,
			uint(0),
			true,
		},
		{
			"false and int16(0) are equivalent",
			false,
			int16(0),
			true,
		},
		{
			"false and float32(0) are not equivalent",
			false,
			float32(0),
			false,
		},
		{
			"false and float64(0) are not equivalent",
			false,
			float64(0),
			false,
		},
		{
			"1 and \"1\" are equivalent",
			1,
			"1",
			true,
		},
		{
			"1 and float32(1.0) are not equivalent",
			1,
			float32(1.0),
			false,
		},
		{
			"uint(1) and float32(1.0) are not equivalent",
			uint(1),
			float32(1.0),
			false,
		},
		{
			"1 and float64(1.0) are not equivalent",
			1,
			float64(1.0),
			false,
		},
		{
			"uint(1) and float64(1.0) are not equivalent",
			uint(1),
			float64(1.0),
			false,
		},
		{
			"\"3.1415\" and float32(3.1415) are equivalent",
			"3.1415",
			float32(3.1415),
			true,
		},
		{
			"\"3.1415\" and float64(3.1415) are equivalent",
			"3.1415",
			float64(3.1415),
			true,
		},
		{
			"[]byte(\"3.1415\") and float32(3.1415) are equivalent",
			[]byte("3.1415"),
			float32(3.1415),
			true,
		},
		{
			"[]byte(\"3.1415\") and float64(3.1415) are equivalent",
			[]byte("3.1415"),
			float64(3.1415),
			true,
		},
		{
			"\"foo\" and []byte(\"foo\") are equivalent",
			"foo",
			[]byte("foo"),
			true,
		},
		{
			"two structs with equivalent fields are not equivalent",
			struct{ a int }{a: 1},
			struct{ a int }{a: 1},
			false,
		},
		{
			"two slices with equivalent elements are not equivalent",
			[]string{"1"},
			[]int{1},
			false,
		},
		{
			"nil and pointer to nil are not equivalent",
			nil,
			sNil,
			false,
		},
		{
			"nil and pointer to non-nil are not equivalent",
			nil,
			sNotNil,
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(c.exp, delta.ScalarEquivalent(c.a, c.b))
			// Test the transitive property of equality
			require.Equal(c.exp, delta.ScalarEquivalent(c.b, c.a))
		})
	}
}

func TestScalar_DecimalPlacesInFloatString(t *testing.T) {
	cases := []struct {
		name      string
		subject   string
		expErr    string
		expPlaces int
	}{
		{
			name:    "invalid float string returns error",
			subject: "not a float",
			expErr:  "invalid syntax",
		},
		{
			name:      "no decimal places returns 0",
			subject:   "12345",
			expPlaces: 0,
		},
		{
			name:      "4 decimal places returns 4",
			subject:   "3.1415",
			expPlaces: 4,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			got, err := delta.DecimalPlacesInFloatString(c.subject)
			if c.expErr != "" {
				require.ErrorContains(err, c.expErr)
			} else {
				require.Nil(err)
				require.Equal(c.expPlaces, got)
			}
		})
	}
}
