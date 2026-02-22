package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityAgentReview(t *testing.T) {
	t.Run("passes when no security issues found", func(t *testing.T) {
		findings := ReviewSecurity(context.Background(), "diff", stubRunner{out: "[]"})
		require.Len(t, findings, 1)
		assert.Equal(t, Pass, findings[0].Verdict)
	})

	t.Run("returns warn finding for SQL injection risk", func(t *testing.T) {
		raw := `[{"agent":"security","verdict":"warn","severity":"medium","finding":"SQL injection risk","file":"db.go:10","fix":"use parameterized query"}]`
		findings := ReviewSecurity(context.Background(), "diff", stubRunner{out: raw})
		require.Len(t, findings, 1)
		assert.Equal(t, Warn, findings[0].Verdict)
	})

	t.Run("returns block finding for hardcoded secret", func(t *testing.T) {
		raw := `[{"agent":"security","verdict":"block","severity":"critical","finding":"hardcoded secret","file":"config.go:3","fix":"use env var"}]`
		findings := ReviewSecurity(context.Background(), "diff", stubRunner{out: raw})
		require.Len(t, findings, 1)
		assert.Equal(t, Block, findings[0].Verdict)
	})

	t.Run("skips gracefully when no runner is configured", func(t *testing.T) {
		findings := ReviewSecurity(context.Background(), "diff", nil)
		require.Len(t, findings, 1)
		assert.Equal(t, Warn, findings[0].Verdict)
		assert.Contains(t, findings[0].Finding, "skipped")
	})

	t.Run("skips gracefully on runner error", func(t *testing.T) {
		findings := ReviewSecurity(context.Background(), "diff", stubRunner{err: errors.New("timeout")})
		require.Len(t, findings, 1)
		assert.Equal(t, Warn, findings[0].Verdict)
		assert.Contains(t, findings[0].Finding, "skipped")
	})

	t.Run("enforces agent name regardless of model response", func(t *testing.T) {
		raw := `[{"agent":"software-engineering","verdict":"warn","severity":"low","finding":"issue","file":"","fix":""}]`
		findings := ReviewSecurity(context.Background(), "diff", stubRunner{out: raw})
		require.Len(t, findings, 1)
		assert.Equal(t, "security", findings[0].Agent)
	})

	t.Run("handles markdown-fenced JSON in model response", func(t *testing.T) {
		raw := "```json\n[]\n```"
		findings := ReviewSecurity(context.Background(), "diff", stubRunner{out: raw})
		require.Len(t, findings, 1)
		assert.Equal(t, Pass, findings[0].Verdict)
	})

	t.Run("skips with warn when model response is unparseable", func(t *testing.T) {
		findings := ReviewSecurity(context.Background(), "diff", stubRunner{out: "not json"})
		require.Len(t, findings, 1)
		assert.Equal(t, Warn, findings[0].Verdict)
		assert.Contains(t, findings[0].Finding, "unparseable")
	})
}
