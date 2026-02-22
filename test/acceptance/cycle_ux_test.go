package acceptance

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCycleIncludesUXAgentWhenDiffHasExportedAPI(t *testing.T) {
	cmd := exec.Command(ensembleBin(t), "cycle")
	cmd.Stdin = strings.NewReader(diffWithExportedFunc())
	cmd.Env = envWithout(os.Environ(), "ANTHROPIC_API_KEY")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "expected exit 0: %s", out)

	findings := parseFindings(t, out)
	hasUXAgent := false
	for _, f := range findings {
		if f["agent"] == "ux-design" {
			hasUXAgent = true
		}
	}
	assert.True(t, hasUXAgent, "no ux-design agent finding in output")
}

func TestCycleOmitsUXAgentWhenDiffHasNoExportedAPI(t *testing.T) {
	cmd := exec.Command(ensembleBin(t), "cycle")
	cmd.Stdin = strings.NewReader(diffWithOnlyUnexportedFuncs())
	cmd.Env = envWithout(os.Environ(), "ANTHROPIC_API_KEY")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "expected exit 0: %s", out)

	findings := parseFindings(t, out)
	for _, f := range findings {
		assert.NotEqual(t, "ux-design", f["agent"], "ux-design should not appear for unexported-only diff")
	}
}

func TestCycleUXAgentRunsOfflineWithWarnFallback(t *testing.T) {
	cmd := exec.Command(ensembleBin(t), "cycle")
	cmd.Stdin = strings.NewReader(diffWithExportedFunc())
	cmd.Env = envWithout(os.Environ(), "ANTHROPIC_API_KEY")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "expected exit 0: %s", out)

	findings := parseFindings(t, out)
	for _, f := range findings {
		if f["agent"] == "ux-design" {
			assert.Equal(t, "warn", f["verdict"])
		}
	}
}

func TestCycleCombinesAllFourAgentsWhenDiffHasExportedAPI(t *testing.T) {
	cmd := exec.Command(ensembleBin(t), "cycle")
	cmd.Stdin = strings.NewReader(diffWithExportedFunc())
	cmd.Env = envWithout(os.Environ(), "ANTHROPIC_API_KEY")
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err, "expected exit 0: %s", out)

	findings := parseFindings(t, out)
	agents := make(map[string]bool)
	for _, f := range findings {
		if a, ok := f["agent"].(string); ok {
			agents[a] = true
		}
	}
	assert.True(t, agents["testing-quality"], "missing testing-quality agent")
	assert.True(t, agents["software-engineering"], "missing software-engineering agent")
	assert.True(t, agents["security"], "missing security agent")
	assert.True(t, agents["ux-design"], "missing ux-design agent")
}

func diffWithExportedFunc() string {
	return `diff --git a/internal/foo/foo.go b/internal/foo/foo.go
--- /dev/null
+++ b/internal/foo/foo.go
@@ -0,0 +1,3 @@
+package foo
+func Add(a, b int) int { return a + b }
diff --git a/internal/foo/foo_test.go b/internal/foo/foo_test.go
--- /dev/null
+++ b/internal/foo/foo_test.go
@@ -0,0 +1,5 @@
+package foo_test
+import "testing"
+func TestAdd(t *testing.T) {
+	if Add(1,2) != 3 { t.Fatal() }
+}
`
}

func diffWithOnlyUnexportedFuncs() string {
	return `diff --git a/internal/bar/bar.go b/internal/bar/bar.go
--- /dev/null
+++ b/internal/bar/bar.go
@@ -0,0 +1,3 @@
+package bar
+func multiply(a, b int) int { return a * b }
diff --git a/internal/bar/bar_test.go b/internal/bar/bar_test.go
--- /dev/null
+++ b/internal/bar/bar_test.go
@@ -0,0 +1,5 @@
+package bar_test
+import "testing"
+func TestMultiply(t *testing.T) {
+	if multiply(2,3) != 6 { t.Fatal() }
+}
`
}
