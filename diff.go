package delta

import (
	"reflect"
	"strconv"

	"github.com/jaypipes/delta/fieldpath"
	"github.com/samber/lo"
)

// Differ allows custom types that know how to determine differences between
// two things to override the more generic comparison functionality in
// Delta.CollectAt().
type Differ interface {
	// Diff should return a Delta representing the set of differences between
	// the thing and the supplied parameter. Should return ErrTypeComparison if
	// the supplied parameter cannot be compared to the thing implementing the
	// Differ interface.
	Diff(any) (Delta, error)
}

// Diff returns a Delta representing the set of differences between two
// supplied things `a` and `b`.
//
// If `a` implements the `Differ` interface, then `Diff` will call
// `a.Diff(b)`.
//
// If `a` does not implement the `Differ` interface and the type of `a` and `b`
// is not comparable, ErrTypeComparison will be returned.
func Diff(a, b any) (Delta, error) {
	differ, ok := a.(Differ)
	if ok {
		return differ.Diff(b)
	}
	d := Delta{}
	err := d.CollectAt(fieldpath.FromString(""), a, b)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// isScalar returns true if the supplied subject is a scalar variable type.
func isScalar(subject any) bool {
	if subject == nil {
		return true
	}
	switch subject.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		string, []byte,
		float64:
		return true
	default:
		return false
	}
}

// isSlice returns true if the supplied subject is a slice of something.
func isSlice(subject any) bool {
	return reflect.TypeOf(subject).Kind() == reflect.Slice
}

// areComparable returns true if the types of the supplied subjects are naively
// comparable.
func areComparable(a, b any) bool {
	if a == nil || b == nil {
		// comparing anything to nil is possible...
		return true
	}
	if isScalar(a) {
		return isScalar(b)
	}
	if isScalar(b) {
		return false
	}
	if isSlice(a) && isSlice(b) {
		return areComparable(
			reflect.TypeOf(a).Elem(),
			reflect.TypeOf(b).Elem(),
		)
	}
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

// equalScalars returns true if the supplied subjects are naively equal. "2" is
// equal to 2. []byte("a") is equal to "a". uint(1) is equal to int(1), etc.
func equalScalars(a, b any) bool {
	switch a := a.(type) {
	case int, int8, int16, int32, int64:
		switch b := b.(type) {
		case int, int8, int16, int32, int64:
			av := toInt64(a)
			bv := toInt64(b)
			return av == bv
		case uint, uint8, uint16, uint32, uint64:
			av := toUint64(a)
			bv := toUint64(b)
			return av == bv
		case string:
			av := toInt64(a)
			bi, err := strconv.Atoi(b)
			if err != nil {
				return false
			}
			return av == int64(bi)
		case []byte:
			av := toInt64(a)
			bi, err := strconv.Atoi(string(b))
			if err != nil {
				return false
			}
			return av == int64(bi)
		}
		return false
	case string:
		switch b := b.(type) {
		case int, int8, int16, int32, int64:
			bi := toInt64(b)
			bs := strconv.FormatInt(bi, 10)
			return a == bs
		case uint, uint8, uint16, uint32, uint64:
			bi := toUint64(b)
			bs := strconv.FormatUint(bi, 10)
			return a == bs
		case string:
			return a == b
		case []byte:
			return a == string(b)
		}
	case []byte:
		switch b := b.(type) {
		case int, int8, int16, int32, int64:
			bi := toInt64(b)
			bs := strconv.FormatInt(bi, 10)
			return string(a) == bs
		case uint, uint8, uint16, uint32, uint64:
			bi := toUint64(b)
			bs := strconv.FormatUint(bi, 10)
			return string(a) == bs
		case string:
			return string(a) == b
		case []byte:
			return string(a) == string(b)
		}
	}
	return false
}

// toUint64 takes an interface and returns a uint64
func toUint64(v any) uint64 {
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
	}
	return 0
}

// toInt64 takes an interface and returns an int64
func toInt64(v any) int64 {
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
	}
	return 0
}

// sliceOfAny returns a []any for the supplied scalar slice-valued thing.
func sliceOfAny(subject any) []any {
	switch subject := subject.(type) {
	case []any:
		return subject
	case []rune:
		return lo.ToAnySlice(subject)
	case []string:
		return lo.ToAnySlice(subject)
	case []int:
		return lo.ToAnySlice(subject)
	case []int8:
		return lo.ToAnySlice(subject)
	case []int16:
		return lo.ToAnySlice(subject)
	case []int64:
		return lo.ToAnySlice(subject)
	case []uint:
		return lo.ToAnySlice(subject)
	case []uint8:
		return lo.ToAnySlice(subject)
	case []uint16:
		return lo.ToAnySlice(subject)
	case []uint32:
		return lo.ToAnySlice(subject)
	case []uint64:
		return lo.ToAnySlice(subject)
	case []float32:
		return lo.ToAnySlice(subject)
	case []float64:
		return lo.ToAnySlice(subject)
	}
	return nil
}

// compareSlices returns two booleans comparing the supplied slices for
// equality and sort order.
func compareSlices(
	a []any,
	b []any,
) (bool, bool) {
	if len(a) != len(b) {
		return false, false
	}
	notInB, notInA := lo.Difference(a, b)
	if len(notInB) > 0 || len(notInA) > 0 {
		return false, false
	}
	return true, sameSortOrder(a, b)
}

// sameSortOrder returns true if the supplied slices are in the same sort
// order.
func sameSortOrder(
	a []any,
	b []any,
) bool {
	if len(a) != len(b) {
		return false
	}
	for x, av := range a {
		bv := b[x]
		if !equalScalars(av, bv) {
			return false
		}
	}
	return true
}
