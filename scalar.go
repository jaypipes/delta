package delta

import (
	"fmt"
	"strconv"
	"strings"
)

// IsScalar returns true if the supplied subject is a scalar variable type.
//
// Important: `nil` is a scalar variable but a reference type with a nil
// underlying value is *not* a scalar variable.
func IsScalar(subject any) bool {
	if subject == nil {
		return true
	}
	switch subject.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		string, []byte,
		float32, float64,
		bool:
		return true
	default:
		return false
	}
}

// ScalarEquivalent returns true if the supplied scalar subjects are naively
// equivalent. Naive equivalence means that:
//
// - "2" is equivalent to 2
// - []byte("a") is equivalent to "a".
// - uint(1) is equivalent to int(1)
// - true is equivalent to "true"
// - true is equivalent to "t"
// - true is equivalent to "TRUE"
// - true is equivalent to "1"
// - true is equivalent to any integer type that is equal to 1
// - false is equivalent to "false"
// - false is equivalent to "f"
// - false is equivalent to "FALSE"
// - false is equivalent to "0"
// - false is equivalent to any integer type that is equal to 0.
// - "3.1415" is equivalent to float32(3.1415) and float64(3.1415)
//
// Comparing integer types to float types always returns false.
//
// Comparing float types to float types always returns false because there is
// no way to determine the number of decimal places of precision to use when
// comparing equivalency. This is in contrast to comparing a string to a float.
// When comparing a string to a float, we can determine the number of decimal
// places in the string representation of the floating point number and convert
// the float to a string at that decimal place precision.
//
// Comparing bool types to float types always returns false.
//
// Note: equivalence is transitive. if a is equivalent to b, then b must be
// equivalent to a.
func ScalarEquivalent(a, b any) bool {
	if !IsScalar(a) || !IsScalar(b) {
		return false
	}
	// Two nils are always equivalent. Note that we deliberately check for
	// whether a and b are the nil type instead of using lo.IsNil() which
	// returns true if a supplied reference type has an underlying value of
	// nil.
	if a == nil && b == nil {
		return true
	}
	switch a := a.(type) {
	case int, int8, int16, int32, int64:
		switch b := b.(type) {
		case int, int8, int16, int32, int64:
			ai := ToInt64(a)
			bi := ToInt64(b)
			return ai == bi
		case uint, uint8, uint16, uint32, uint64:
			ai := ToUint64(a)
			bi := ToUint64(b)
			return ai == bi
		case string:
			ai := ToInt64(a)
			bi, err := strconv.Atoi(b)
			if err != nil {
				return false
			}
			return ai == int64(bi)
		case []byte:
			ai := ToInt64(a)
			bi, err := strconv.Atoi(string(b))
			if err != nil {
				return false
			}
			return ai == int64(bi)
		case bool:
			ai := ToInt64(a)
			if ai != int64(0) && ai != int64(1) {
				return false
			}
			return (ai != int64(0)) == b
		}
		return false
	case uint, uint8, uint16, uint32, uint64:
		switch b := b.(type) {
		case int, int8, int16, int32, int64:
			av := ToUint64(a)
			bv := ToUint64(b)
			return av == bv
		case uint, uint8, uint16, uint32, uint64:
			av := ToUint64(a)
			bv := ToUint64(b)
			return av == bv
		case string:
			av := ToInt64(a)
			bi, err := strconv.Atoi(b)
			if err != nil {
				return false
			}
			return av == int64(bi)
		case []byte:
			av := ToInt64(a)
			bi, err := strconv.Atoi(string(b))
			if err != nil {
				return false
			}
			return av == int64(bi)
		case bool:
			ai := ToUint64(a)
			if ai != uint64(0) && ai != uint64(1) {
				return false
			}
			return (ai != uint64(0)) == b
		}
		return false
	case float32:
		switch b := b.(type) {
		case string:
			places, err := DecimalPlacesInFloatString(b)
			if err != nil {
				return false
			}
			strFmt := fmt.Sprintf("%%.%df", places)
			return fmt.Sprintf(strFmt, a) == b
		case []byte:
			places, err := DecimalPlacesInFloatString(string(b))
			if err != nil {
				return false
			}
			strFmt := fmt.Sprintf("%%.%df", places)
			return fmt.Sprintf(strFmt, a) == string(b)
		}
		return false
	case float64:
		switch b := b.(type) {
		case string:
			places, err := DecimalPlacesInFloatString(b)
			if err != nil {
				return false
			}
			strFmt := fmt.Sprintf("%%.%df", places)
			return fmt.Sprintf(strFmt, a) == b
		case []byte:
			places, err := DecimalPlacesInFloatString(string(b))
			if err != nil {
				return false
			}
			strFmt := fmt.Sprintf("%%.%df", places)
			return fmt.Sprintf(strFmt, a) == string(b)
		}
		return false
	case string:
		switch b := b.(type) {
		case int, int8, int16, int32, int64:
			bi := ToInt64(b)
			bs := strconv.FormatInt(bi, 10)
			return a == bs
		case uint, uint8, uint16, uint32, uint64:
			bi := ToUint64(b)
			bs := strconv.FormatUint(bi, 10)
			return a == bs
		case float32:
			places, err := DecimalPlacesInFloatString(a)
			if err != nil {
				return false
			}
			strFmt := fmt.Sprintf("%%.%df", places)
			return fmt.Sprintf(strFmt, b) == a
		case float64:
			places, err := DecimalPlacesInFloatString(a)
			if err != nil {
				return false
			}
			strFmt := fmt.Sprintf("%%.%df", places)
			return fmt.Sprintf(strFmt, b) == a
		case string:
			return a == b
		case []byte:
			return a == string(b)
		case bool:
			ab, err := strconv.ParseBool(a)
			if err != nil {
				return false
			}
			return ab == b
		}
	case []byte:
		switch b := b.(type) {
		case int, int8, int16, int32, int64:
			bi := ToInt64(b)
			bs := strconv.FormatInt(bi, 10)
			return string(a) == bs
		case uint, uint8, uint16, uint32, uint64:
			bi := ToUint64(b)
			bs := strconv.FormatUint(bi, 10)
			return string(a) == bs
		case float32:
			places, err := DecimalPlacesInFloatString(string(a))
			if err != nil {
				return false
			}
			strFmt := fmt.Sprintf("%%.%df", places)
			return fmt.Sprintf(strFmt, b) == string(a)
		case float64:
			places, err := DecimalPlacesInFloatString(string(a))
			if err != nil {
				return false
			}
			strFmt := fmt.Sprintf("%%.%df", places)
			return fmt.Sprintf(strFmt, b) == string(a)
		case string:
			return string(a) == b
		case []byte:
			return string(a) == string(b)
		case bool:
			ab, err := strconv.ParseBool(string(a))
			if err != nil {
				return false
			}
			return ab == b
		}
	case bool:
		switch b := b.(type) {
		case int, int8, int16, int32, int64:
			bi := ToInt64(b)
			if bi != int64(0) && bi != int64(1) {
				return false
			}
			return a == (bi != 0)
		case uint, uint8, uint16, uint32, uint64:
			bi := ToUint64(b)
			if bi != uint64(0) && bi != uint64(1) {
				return false
			}
			return a == (bi != uint64(0))
		case string:
			bb, err := strconv.ParseBool(b)
			if err != nil {
				return false
			}
			return a == bb
		case []byte:
			bb, err := strconv.ParseBool(string(b))
			if err != nil {
				return false
			}
			return a == bb
		case bool:
			return a == b
		}
	}
	return false
}

// ToUint64 takes an interface and returns a uint64
func ToUint64(v any) uint64 {
	switch v := v.(type) {
	case uint64:
		return v
	case uint8:
		return uint64(v)
	case uint16:
		return uint64(v)
	case uint32:
		return uint64(v)
	case uint:
		return uint64(v)
	case float32:
		return uint64(v)
	case float64:
		return uint64(v)
	}
	return 0
}

// ToInt64 takes an interface and returns an int64
func ToInt64(v any) int64 {
	switch v := v.(type) {
	case int64:
		return v
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int:
		return int64(v)
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	}
	return 0
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
