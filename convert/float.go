package convert

import (
	"fmt"
	"strconv"
	"strings"
)

// ToFloatString returns the supplied float64 formatted as a string with the
// required decimal places.
func ToFloatString(f float64, places int) string {
	strFmt := fmt.Sprintf("%%.%df", places)
	return fmt.Sprintf(strFmt, f)
}

// DecimalPlacesInFloatString returns the number of decimal places in the
// supplied string representing a floating point number.
func DecimalPlacesInFloatString(subject string) (int, error) {
	_, err := strconv.ParseFloat(subject, 64)
	if err != nil {
		return -1, err
	}
	idx := strings.Index(subject, ".")
	if idx != -1 {
		return len(subject) - idx - 1, nil
	}
	return 0, nil
}
