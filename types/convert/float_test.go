package convert_test

import (
	"testing"

	"github.com/jaypipes/delta/types/convert"
	"github.com/stretchr/testify/assert"
)

func TestDecimalPlacesInFloatString(t *testing.T) {
	cases := []struct {
		subject   string
		expErr    string
		expPlaces int
	}{
		{
			subject: "not a float",
			expErr:  "invalid syntax",
		},
		{
			subject:   "12345",
			expPlaces: 0,
		},
		{
			subject:   "3.1415",
			expPlaces: 4,
		},
	}
	assert := assert.New(t)
	for _, c := range cases {
		got, err := convert.DecimalPlacesInFloatString(c.subject)
		if c.expErr != "" {
			assert.ErrorContains(err, c.expErr)
		} else {
			assert.Nil(err)
			assert.Equal(c.expPlaces, got)
		}
	}
}
