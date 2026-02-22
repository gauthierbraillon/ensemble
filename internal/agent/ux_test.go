package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUXDesignAgentReview(t *testing.T) {
	t.Run("silent when diff contains no exported API change", func(t *testing.T) {
		diff := `diff --git a/internal/foo/foo.go b/internal/foo/foo.go
+++ b/internal/foo/foo.go
+func add(a, b int) int { return a + b }`
		findings := ReviewUX(context.Background(), diff, stubRunner{out: "[]"})
		assert.Empty(t, findings)
	})

	t.Run("passes when no API design issues found", func(t *testing.T) {
		findings := ReviewUX(context.Background(), diffWithExportedFuncUnit(), stubRunner{out: "[]"})
		require.Len(t, findings, 1)
		assert.Equal(t, Pass, findings[0].Verdict)
	})

	t.Run("returns warn finding for unclear exported name", func(t *testing.T) {
		raw := `[{"agent":"ux-design","verdict":"warn","severity":"low","finding":"unclear name","file":"api.go:1","fix":"rename to GetUser"}]`
		findings := ReviewUX(context.Background(), diffWithExportedFuncUnit(), stubRunner{out: raw})
		require.Len(t, findings, 1)
		assert.Equal(t, Warn, findings[0].Verdict)
	})

	t.Run("returns block finding for API convention violation", func(t *testing.T) {
		raw := `[{"agent":"ux-design","verdict":"block","severity":"high","finding":"API breaks convention","file":"api.go:5","fix":"follow REST naming"}]`
		findings := ReviewUX(context.Background(), diffWithExportedFuncUnit(), stubRunner{out: raw})
		require.Len(t, findings, 1)
		assert.Equal(t, Block, findings[0].Verdict)
	})

	t.Run("skips gracefully when no runner is configured", func(t *testing.T) {
		findings := ReviewUX(context.Background(), diffWithExportedFuncUnit(), nil)
		require.Len(t, findings, 1)
		assert.Equal(t, Warn, findings[0].Verdict)
		assert.Contains(t, findings[0].Finding, "skipped")
	})

	t.Run("skips gracefully on runner error", func(t *testing.T) {
		findings := ReviewUX(context.Background(), diffWithExportedFuncUnit(), stubRunner{err: errors.New("timeout")})
		require.Len(t, findings, 1)
		assert.Equal(t, Warn, findings[0].Verdict)
		assert.Contains(t, findings[0].Finding, "skipped")
	})

	t.Run("enforces agent name regardless of model response", func(t *testing.T) {
		raw := `[{"agent":"software-engineering","verdict":"warn","severity":"low","finding":"issue","file":"","fix":""}]`
		findings := ReviewUX(context.Background(), diffWithExportedFuncUnit(), stubRunner{out: raw})
		require.Len(t, findings, 1)
		assert.Equal(t, "ux-design", findings[0].Agent)
	})

	t.Run("handles markdown-fenced JSON in model response", func(t *testing.T) {
		raw := "```json\n[]\n```"
		findings := ReviewUX(context.Background(), diffWithExportedFuncUnit(), stubRunner{out: raw})
		require.Len(t, findings, 1)
		assert.Equal(t, Pass, findings[0].Verdict)
	})
}

func TestExportedAPIChangeDetection(t *testing.T) {
	t.Run("detects added exported function", func(t *testing.T) {
		assert.True(t, hasExportedAPIChange("+func Add(a, b int) int { return a + b }"))
	})

	t.Run("ignores added unexported function", func(t *testing.T) {
		assert.False(t, hasExportedAPIChange("+func add(a, b int) int { return a + b }"))
	})

	t.Run("ignores exported symbols in test files", func(t *testing.T) {
		diff := `+++ b/foo_test.go
+func TestAdd(t *testing.T) {}
+func Add(a, b int) int { return a + b }`
		assert.False(t, hasExportedAPIChange(diff))
	})

	t.Run("detects added exported type", func(t *testing.T) {
		assert.True(t, hasExportedAPIChange("+type Config struct { Host string }"))
	})

	t.Run("detects added exported constant", func(t *testing.T) {
		assert.True(t, hasExportedAPIChange("+const MaxRetries = 3"))
	})

	t.Run("ignores removed exported symbols — deletions do not trigger review", func(t *testing.T) {
		assert.False(t, hasExportedAPIChange("-func Remove() {}"))
	})
}

func diffWithExportedFuncUnit() string {
	return `diff --git a/internal/foo/foo.go b/internal/foo/foo.go
+++ b/internal/foo/foo.go
+func Add(a, b int) int { return a + b }`
}
