package fieldpath

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

var (
	ErrInvalid = fmt.Errorf("invalid field path")
)

// FromAny returns a new FieldPath from a dotted-notation string, e.g.
// "spec.generation", or a slice of strings, e.g. []string{"spec",
// "generation"} or a FieldPath object.
func FromAny(subject any) (FieldPath, error) {
	switch subject := subject.(type) {
	case []string:
		return FieldPath(subject), nil
	case FieldPath:
		return subject, nil
	case string:
		return FieldPath(strings.Split(subject, ".")), nil
	default:
		return nil, fmt.Errorf(
			"%w: unsupported type %T", ErrInvalid, subject,
		)
	}
}

// FromString returns a new FieldPath from a dotted-notation string, e.g.
// "spec.generation".
func FromString(subject string) FieldPath {
	return FieldPath(strings.Split(subject, "."))
}

// Prefixed returns a new FieldPath prefixed with the supplied prefix.
//
// For example, assume a FieldPath `fp` of `[]string{"generation"}`, calling
// `fieldpath.Prefixed(fp, "spec") would return a new FieldPath of
// `[]string{"spec", "generation"}`.
func Prefixed(subject FieldPath, prefixes ...string) FieldPath {
	prefixes = lo.Compact(prefixes) // removes empty strings...
	fp := FieldPath(prefixes)
	fp = append(fp, subject...)
	return fp
}
