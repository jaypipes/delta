package diff

import (
	"strings"

	"github.com/samber/lo"

	"github.com/jaypipes/delta/diff/fieldpath"
)

// ZeroGen is a sentinel type used in creating Deltas for new things.
type ZeroGen struct{}

var (
	Zero = ZeroGen{}
)

// Delta represents the difference between two comparable things.
type Delta []Difference

// Push adds a single Difference to the Delta.
func (d *Delta) Push(diff Difference) {
	*d = append(*d, diff)
}

// ExtendAt adds the supplied Delta's differences at a supplied FieldPath.
func (d *Delta) ExtendAt(
	at fieldpath.FieldPath,
	other Delta,
) {
	for _, diff := range other.Differences() {
		*d = append(
			*d,
			Difference{
				FieldPath: fieldpath.Prefixed(diff.FieldPath, at...),
				Kind:      diff.Kind,
				From:      diff.From,
				To:        diff.To,
			},
		)
	}
}

// Differences returns the Delta's slice of Differences
func (d Delta) Differences() []Difference {
	if d == nil {
		return []Difference{}
	}
	return d
}

// DifferencesAt returns the Differences in the Delta at the supplied path, or
// nil if there was no Difference at the supplied path.
func (d Delta) DifferencesAt(
	subject fieldpath.FieldPath,
) []Difference {
	return lo.Filter(d, func(diff Difference, _ int) bool {
		return diff.FieldPath.Contains(subject)
	})
}

// DifferencesAtAny returns the Differences in the Delta at any of the supplied
// paths, or nil if there was no Difference at the supplied path.
func (d Delta) DifferencesAtAny(
	subjects ...fieldpath.FieldPath,
) []Difference {
	return lo.Filter(d, func(diff Difference, _ int) bool {
		return diff.FieldPath.ContainsAny(subjects)
	})
}

// DifferencesExcept returns the Differences in the Delta other than any of the
// supplied paths.
func (d Delta) DifferencesExcept(
	exceptPaths ...fieldpath.FieldPath,
) []Difference {
	return lo.Filter(d, func(diff Difference, _ int) bool {
		return !diff.FieldPath.ContainsAny(exceptPaths)
	})
}

// Empty returns true if the Delta has no differences.
func (d Delta) Empty() bool {
	return len(d) == 0
}

// Different returns true if there are any differences in the Delta.
func (d Delta) Different() bool {
	return !d.Empty()
}

// DifferentAt returns whether there are differences at the supplied path in
// the things under comparison.
func (d Delta) DifferentAt(subject fieldpath.FieldPath) bool {
	if d.Empty() {
		return false
	}
	for _, diff := range d {
		if diff.FieldPath.Contains(subject) {
			return true
		}
	}
	return false
}

// DifferentAt returns whether there are differences at any of the the supplied
// paths in the things under comparison.
func (d Delta) DifferentAtAny(subjects ...fieldpath.FieldPath) bool {
	if d.Empty() {
		return false
	}
	for _, diff := range d {
		if diff.FieldPath.ContainsAny(subjects) {
			return true
		}
	}
	return false
}

// DifferentExcept returns true if there are differences in the things under
// comparison at any path *except* the supplied paths.
//
// This method is useful when comparing different objects where you only care
// about whether certain fields (or sets of fields) are different but not
// others. For example, assume some object with the following pseudo-schema:
//
// | {
// |   "spec": {
// |     "replicas": <int>,
// |     "generation": <int>
// |   },
// |   "status": {
// |     "readyReplicas": <int>,
// |     "lastHeartbeat": <timestamp>
// |   },
// | }
//
// If we wanted to know whether two objects with the above schema had different
// spec.replicas or spec.generation field values but don't care if the
// status.lastHeartbeat or status.readyReplicas field values are different, we
// could call `Delta.DifferentExcept("status") to determine that.
func (d Delta) DifferentExcept(
	exceptPaths ...fieldpath.FieldPath,
) bool {
	numDiffs := len(d)
	if numDiffs == 0 {
		return false
	} else if numDiffs > len(exceptPaths) {
		return true
	}
	foundExcepts := 0
	for _, diff := range d {
		for _, exceptPath := range exceptPaths {
			if diff.FieldPath.Contains(exceptPath) {
				foundExcepts++
			}
		}
	}
	return foundExcepts != numDiffs
}

// String returns a nicely-formatted string showing the differences contained
// in the Delta.
func (d Delta) String() string {
	if d.Empty() {
		return "<empty>"
	}
	b := strings.Builder{}
	for x, diff := range d {
		if x > 0 {
			b.WriteRune('\n')
		}
		b.WriteString(diff.String())
	}
	return b.String()
}
