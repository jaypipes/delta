package delta

import (
	"fmt"

	"github.com/jaypipes/delta/fieldpath"
)

// DifferenceType differentiates kinds of Differences
type DifferenceType int

const (
	// DifferenceTypeAdd represents an addition of some field or item in a
	// list/map field value type.
	DifferenceTypeAdd DifferenceType = iota
	// DifferenceTypeRemove represents the removal of some field or item in a
	// list/map field value type.
	DifferenceTypeRemove
	// DifferenceTypeModify represents a modification of some field or item in
	// a list/map field value type.
	DifferenceTypeModify
	// DifferenceTypeOrder represents a change in the order of elements in a
	// sorted list field.
	DifferenceTypeOrder
)

func (t DifferenceType) String() string {
	switch t {
	case DifferenceTypeAdd:
		return "ADD"
	case DifferenceTypeRemove:
		return "REMOVE"
	case DifferenceTypeModify:
		return "MODIFY"
	case DifferenceTypeOrder:
		return "ORDER"
	}
	return ""
}

// Difference represents a single difference between two things.
type Difference struct {
	// FieldPath contains the FieldPath where this difference exists.
	FieldPath fieldpath.FieldPath
	// Type contains the type of difference.
	Type DifferenceType
	// From contains the "A" or "from" object's value at the FieldPath.
	From any
	// To contains the "B" or "to" object's value at the FieldPath.
	To any
}

func (d Difference) String() string {
	fp := d.FieldPath.String()
	if len(fp) > 0 {
		fp += " "
	}
	switch d.Type {
	case DifferenceTypeAdd:
		return fmt.Sprintf("%sADD new value: %v", fp, d.To)
	case DifferenceTypeRemove:
		return fmt.Sprintf("%sREMOVE previous value: %v", fp, d.From)
	case DifferenceTypeModify:
		return fmt.Sprintf("%sMODIFY %v -> %v", fp, d.From, d.To)
	case DifferenceTypeOrder:
		return fmt.Sprintf("%sSORT ORDER", fp)
	}
	return ""
}
