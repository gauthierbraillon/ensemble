package acceptance

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ensembleBin(t *testing.T) string {
	t.Helper()
	bin := "../../bin/ensemble"
	if _, err := os.Stat(bin); os.IsNotExist(err) {
		out, err := exec.Command("go", "build", "-o", bin, "../../.").CombinedOutput()
		require.NoError(t, err, "build failed: %s", out)
	}
	return bin
}

func parseFindings(t *testing.T, out []byte) []map[string]interface{} {
	t.Helper()
	var findings []map[string]interface{}
	for _, line := range bytes.Split(bytes.TrimSpace(out), []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		var f map[string]interface{}
		if err := json.Unmarshal(line, &f); err == nil {
			findings = append(findings, f)
		}
	}
	return findings
}

func envWithout(env []string, key string) []string {
	var out []string
	prefix := key + "="
	for _, kv := range env {
		if !strings.HasPrefix(kv, prefix) {
			out = append(out, kv)
		}
	}
	return out
}

func diffWithTest() string {
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

func diffWithoutTest() string {
	return `diff --git a/internal/bar/bar.go b/internal/bar/bar.go
--- /dev/null
+++ b/internal/bar/bar.go
@@ -0,0 +1,3 @@
+package bar
+func Multiply(a, b int) int { return a * b }
`
}

func TestCycleCommand(t *testing.T) {
	t.Run("each finding is one JSON line for CI integration", func(t *testing.T) {
		cmd := exec.Command(ensembleBin(t), "cycle")
		cmd.Stdin = strings.NewReader("diff --git a/x.go b/x.go\n")
		out, _ := cmd.CombinedOutput()

		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if line == "" {
				continue
			}
			var m map[string]interface{}
			assert.NoError(t, json.Unmarshal([]byte(line), &m), "not valid JSON: %q", line)
		}
	})

	t.Run("passes when every implementation file has a matching test", func(t *testing.T) {
		cmd := exec.Command(ensembleBin(t), "cycle")
		cmd.Stdin = strings.NewReader(diffWithTest())
		out, err := cmd.CombinedOutput()
		assert.NoError(t, err, "unexpected block: %s", out)

		for _, f := range parseFindings(t, out) {
			assert.NotEqual(t, "block", f["verdict"])
		}
	})

	t.Run("blocks when an implementation file has no test", func(t *testing.T) {
		cmd := exec.Command(ensembleBin(t), "cycle")
		cmd.Stdin = strings.NewReader(diffWithoutTest())
		out, _ := cmd.CombinedOutput()

		assert.Equal(t, 1, cmd.ProcessState.ExitCode(), "expected exit 1: %s", out)

		hasBlock := false
		for _, f := range parseFindings(t, out) {
			if f["verdict"] == "block" {
				hasBlock = true
			}
		}
		assert.True(t, hasBlock)
	})

	t.Run("runs offline with warn fallback when no API key is configured", func(t *testing.T) {
		cmd := exec.Command(ensembleBin(t), "cycle")
		cmd.Stdin = strings.NewReader(diffWithTest())
		cmd.Env = envWithout(os.Environ(), "ANTHROPIC_API_KEY")
		out, err := cmd.CombinedOutput()
		assert.NoError(t, err, "expected exit 0: %s", out)

		findings := parseFindings(t, out)
		var sweFindings []map[string]interface{}
		for _, f := range findings {
			if f["agent"] == "software-engineering" {
				sweFindings = append(sweFindings, f)
			}
		}
		require.Len(t, sweFindings, 1)
		assert.Equal(t, "warn", sweFindings[0]["verdict"])
	})

	t.Run("combines all active agents in a single output stream", func(t *testing.T) {
		cmd := exec.Command(ensembleBin(t), "cycle")
		cmd.Stdin = strings.NewReader(diffWithTest())
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
	})
}
