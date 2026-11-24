package convert_test

import (
	"testing"

	"github.com/jaypipes/delta/convert"
	"github.com/stretchr/testify/require"
)

func TestDecimalPlacesInFloatString(t *testing.T) {
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
			got, err := convert.DecimalPlacesInFloatString(c.subject)
			if c.expErr != "" {
				require.ErrorContains(err, c.expErr)
			} else {
				require.Nil(err)
				require.Equal(c.expPlaces, got)
			}
		})
	}
}
