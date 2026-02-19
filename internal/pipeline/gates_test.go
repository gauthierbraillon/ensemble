package pipeline_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gauthierbraillon/ensemble/internal/pipeline"
)

func TestGatesReturnsNineGates(t *testing.T) {
	gates := pipeline.Gates()
	require.Len(t, gates, 9, "pipeline must have exactly 9 quality gates")
}

func TestEachGateHasRequiredFields(t *testing.T) {
	for _, gate := range pipeline.Gates() {
		assert.NotEmpty(t, gate.Name, "gate must have a name")
		assert.NotEmpty(t, gate.MakeTarget, "gate must have a make target")
		assert.NotEmpty(t, gate.Description, "gate must have a description")
	}
}

func TestGateOrderFastFailFirst(t *testing.T) {
	gates := pipeline.Gates()
	// Lint and typecheck must come before test and vulncheck.
	// Cheap static checks before expensive runtime checks.
	indexOf := func(target string) int {
		for i, g := range gates {
			if g.MakeTarget == target {
				return i
			}
		}
		return -1
	}

	assert.Less(t, indexOf("lint"), indexOf("test"), "lint must run before tests")
	assert.Less(t, indexOf("typecheck"), indexOf("test"), "typecheck must run before tests")
	assert.Less(t, indexOf("secrets"), indexOf("build"), "secret scan must run before build")
	assert.Less(t, indexOf("sast"), indexOf("build"), "SAST must run before build")
	assert.Less(t, indexOf("build"), indexOf("vulncheck"), "build must precede vuln check")
}
