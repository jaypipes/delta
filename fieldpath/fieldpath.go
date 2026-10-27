package fieldpath

import "strings"

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
