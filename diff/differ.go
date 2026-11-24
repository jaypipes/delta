package diff

import (
	"github.com/jaypipes/delta/cmp"
)

// Differ allows custom types that know how to determine differences between
// two things to override the more generic comparison functionality in
// Delta.CollectAt().
type Differ interface {
	// Diff should return a Delta representing the set of differences between
	// the thing and the supplied parameter. Should return ErrTypeComparison if
	// the supplied parameter cannot be compared to the thing implementing the
	// Differ interface.
	Diff(any, ...cmp.Option) (Delta, error)
}
