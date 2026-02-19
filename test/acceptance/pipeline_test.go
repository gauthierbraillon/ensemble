// Package acceptance tests verify system-level behaviour.
// This file documents and enforces the CD pipeline contract:
// all Minimum CD required gates must exist in both the Makefile and CI workflow.
package acceptance

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// requiredMakeTargets are the 9 quality gates every commit must pass.
// Their order in the Makefile matters: gates run sequentially, fast-fail first.
var requiredMakeTargets = []string{
	"lint",
	"typecheck",
	"secrets",
	"sast",
	"build",
	"test",
	"vulncheck",
	"test-contracts",
	"test-schema",
	"ci", // orchestrates all gates in order
}

// requiredCISteps must appear in .github/workflows/ci.yml.
// Each step name maps directly to a Makefile target.
var requiredCISteps = []string{
	"lint",
	"typecheck",
	"secrets",
	"sast",
	"build",
	"test",
	"vulncheck",
	"test-contracts",
	"test-schema",
}

func TestMakefileHasAllQualityGates(t *testing.T) {
	content, err := os.ReadFile("../../Makefile")
	require.NoError(t, err, "Makefile must exist")

	makefile := string(content)
	for _, target := range requiredMakeTargets {
		require.Contains(t, makefile, target+":", "Makefile must define target: "+target)
	}
}

func TestCIWorkflowHasAllQualityGates(t *testing.T) {
	content, err := os.ReadFile("../../.github/workflows/ci.yml")
	require.NoError(t, err, ".github/workflows/ci.yml must exist")

	var workflow map[string]interface{}
	require.NoError(t, yaml.Unmarshal(content, &workflow), "ci.yml must be valid YAML")

	workflowStr := string(content)
	for _, step := range requiredCISteps {
		require.True(t,
			strings.Contains(workflowStr, "make "+step),
			"ci.yml must invoke make %s", step,
		)
	}
}

func TestCIWorkflowTriggersOnMainPush(t *testing.T) {
	content, err := os.ReadFile("../../.github/workflows/ci.yml")
	require.NoError(t, err, ".github/workflows/ci.yml must exist")

	var workflow map[string]interface{}
	require.NoError(t, yaml.Unmarshal(content, &workflow), "ci.yml must be valid YAML")

	on, ok := workflow["on"]
	require.True(t, ok, "workflow must have 'on' trigger")

	onStr := fmt.Sprint(on)
	require.Contains(t, onStr, "push", "workflow must trigger on push")
}
