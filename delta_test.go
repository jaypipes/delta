package delta_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jaypipes/delta"
	"github.com/jaypipes/delta/fieldpath"
)

var (
	fpSpecReplicas        = fieldpath.FieldPath{"spec", "replicas"}
	fpSpecGeneration      = fieldpath.FieldPath{"spec", "generation"}
	fpStatusReadyReplicas = fieldpath.FieldPath{"status", "readyReplicas"}
)

var (
	diffAB_SpecReplicas = delta.Difference{
		FieldPath: fpSpecReplicas,
		Type:      delta.DifferenceTypeModify,
		From:      1,
		To:        2,
	}
	diffAB_StatusReadyReplicas = delta.Difference{
		FieldPath: fpStatusReadyReplicas,
		Type:      delta.DifferenceTypeModify,
		From:      0,
		To:        1,
	}
	diffAC_StatusReadyReplicas = delta.Difference{
		FieldPath: fpStatusReadyReplicas,
		Type:      delta.DifferenceTypeModify,
		From:      0,
		To:        1,
	}
	diffBC_SpecReplicas = delta.Difference{
		FieldPath: fpSpecReplicas,
		Type:      delta.DifferenceTypeModify,
		From:      2,
		To:        1,
	}
	diffBC_StatusReadyReplicas = delta.Difference{
		FieldPath: fpStatusReadyReplicas,
		Type:      delta.DifferenceTypeModify,
		From:      1,
		To:        0,
	}
)

var (
	emptyDelta = delta.Delta{}
	// B is different from A in spec.replicas and status.readyReplicas
	deltaAB = delta.Delta{
		diffAB_SpecReplicas,
		diffAB_StatusReadyReplicas,
	}
	// C is different from A only in status.readyReplicas
	deltaAC = delta.Delta{
		diffAC_StatusReadyReplicas,
	}
	// C is different from B in spec.replicas and status.readyReplicas
	deltaBC = delta.Delta{
		diffBC_SpecReplicas,
		diffBC_StatusReadyReplicas,
	}
)

func TestDelta_Different(t *testing.T) {
	cases := []struct {
		name    string
		subject delta.Delta
		exp     bool
	}{
		{
			"empty delta",
			emptyDelta,
			false,
		},
		{
			"non-empty delta",
			deltaAB,
			true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(c.exp, c.subject.Different())
		})
	}
}

func TestDelta_ExtendAt(t *testing.T) {
	require := require.New(t)
	abc := delta.Delta{}
	abc.ExtendAt(fieldpath.FromString("ab"), deltaAB)
	abc.ExtendAt(fieldpath.FromString("ac"), deltaAC)
	expected := delta.Delta{
		{
			FieldPath: fieldpath.FromString("ab.spec.replicas"),
			Type:      delta.DifferenceTypeModify,
			From:      1,
			To:        2,
		},
		{
			FieldPath: fieldpath.FromString("ab.status.readyReplicas"),
			Type:      delta.DifferenceTypeModify,
			From:      0,
			To:        1,
		},
		{
			FieldPath: fieldpath.FromString("ac.status.readyReplicas"),
			Type:      delta.DifferenceTypeModify,
			From:      0,
			To:        1,
		},
	}
	require.Equal(expected, abc)
}

func TestDelta_DifferentAt(t *testing.T) {
	cases := []struct {
		name    string
		subject delta.Delta
		path    fieldpath.FieldPath
		exp     bool
	}{
		{
			"empty delta",
			emptyDelta,
			fpSpecReplicas,
			false,
		},
		{
			"AB spec.replicas",
			deltaAB,
			fpSpecReplicas,
			true,
		},
		{
			"AB spec.generation",
			deltaAB,
			fpSpecGeneration,
			false,
		},
		{
			"AB status.readyReplicas",
			deltaAB,
			fpStatusReadyReplicas,
			true,
		},
		{
			"AC spec.replicas",
			deltaAC,
			fpSpecReplicas,
			false,
		},
		{
			"AC spec.generation",
			deltaAC,
			fpSpecGeneration,
			false,
		},
		{
			"AC status.readyReplicas",
			deltaAC,
			fpStatusReadyReplicas,
			true,
		},
		{
			"BC spec.replicas",
			deltaBC,
			fpSpecReplicas,
			true,
		},
		{
			"BC spec.generation",
			deltaBC,
			fpSpecGeneration,
			false,
		},
		{
			"BC status.readyReplicas",
			deltaBC,
			fpStatusReadyReplicas,
			true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(c.exp, c.subject.DifferentAt(c.path))
		})
	}
}

func TestDelta_DifferentExcept(t *testing.T) {
	cases := []struct {
		name    string
		subject delta.Delta
		paths   []fieldpath.FieldPath
		exp     bool
	}{
		{
			"empty delta",
			emptyDelta,
			[]fieldpath.FieldPath{fpSpecReplicas},
			false,
		},
		{
			"AB except spec.replicas",
			deltaAB,
			[]fieldpath.FieldPath{fpSpecReplicas},
			true, // because also diff at status.readyReplicas
		},
		{
			"AB except spec.replicas and status.readyReplicas",
			deltaAB,
			[]fieldpath.FieldPath{fpSpecReplicas, fpStatusReadyReplicas},
			false, // because above paths are the only diffs
		},
		{
			"AC except status.readyReplicas",
			deltaAC,
			[]fieldpath.FieldPath{fpStatusReadyReplicas},
			false, // because above path is the only diffs
		},
		{
			"no supplied field paths and non-empty delta",
			deltaAB,
			nil,
			true, // because there are differences
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(c.exp, c.subject.DifferentExcept(c.paths...))
		})
	}
}

func TestDelta_DifferencesExcept(t *testing.T) {
	cases := []struct {
		name    string
		subject delta.Delta
		paths   []fieldpath.FieldPath
		exp     []delta.Difference
	}{
		{
			"empty delta",
			emptyDelta,
			[]fieldpath.FieldPath{fpSpecReplicas},
			[]delta.Difference{},
		},
		{
			"AB except spec.replicas",
			deltaAB,
			[]fieldpath.FieldPath{fpSpecReplicas},
			[]delta.Difference{
				{
					FieldPath: fieldpath.FromString("status.readyReplicas"),
					Type:      delta.DifferenceTypeModify,
					From:      0,
					To:        1,
				},
			},
		},
		{
			"AB except spec.replicas and status.readyReplicas",
			deltaAB,
			[]fieldpath.FieldPath{fpSpecReplicas, fpStatusReadyReplicas},
			[]delta.Difference{},
		},
		{
			"AC except status.readyReplicas",
			deltaAC,
			[]fieldpath.FieldPath{fpStatusReadyReplicas},
			[]delta.Difference{}, // because above path is the only diffs
		},
		{
			"no supplied field paths and non-empty delta",
			deltaAB,
			nil,
			[]delta.Difference{
				{
					FieldPath: fieldpath.FromString("spec.replicas"),
					Type:      delta.DifferenceTypeModify,
					From:      1,
					To:        2,
				},
				{
					FieldPath: fieldpath.FromString("status.readyReplicas"),
					Type:      delta.DifferenceTypeModify,
					From:      0,
					To:        1,
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(c.exp, c.subject.DifferencesExcept(c.paths...))
		})
	}
}
