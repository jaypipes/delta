package fieldpath

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

var (
	ErrInvalid = fmt.Errorf("invalid field path")
)

// FieldPath provides a "route" to a particular field within a compared thing.
type FieldPath []string

// String returns the FieldPath as a dotted-notation string.
func (p FieldPath) String() string {
	return strings.Join(p, ".")
}

// Push adds a new part to the FieldPath.
func (p *FieldPath) Push(part string) {
	*p = append(*p, part)
}

// Pop removes the last part from the FieldPath
func (p *FieldPath) Pop() {
	fp := *p
	if len(fp) > 0 {
		*p = fp[:len(fp)-1]
	}
}

// Contains returns true if the supplied FieldPath matches parts of this
// FieldPath.
//
//	e.g. if the FieldPath p represents "A.B":
//		subject "A" -> true
//		subject "A.B" -> true
//		subject "A.B.C" -> false
//		subject "B" -> false
//		subject "A.C" -> false
func (p FieldPath) Contains(subject FieldPath) bool {
	if len(subject) > len(p) {
		return false
	}
	for i, s := range subject {
		if p[i] != s {
			return false
		}
	}
	return true
}

// ContainsAny returns true if any of the supplied FieldPaths matches parts of
// this FieldPath.
func (p FieldPath) ContainsAny(subjects []FieldPath) bool {
	for _, subject := range subjects {
		if p.Contains(subject) {
			return true
		}
	}
	return false
}

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
