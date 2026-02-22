package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubRunner struct {
	out string
	err error
}

func (s stubRunner) Run(_ context.Context, _ string) (string, error) {
	return s.out, s.err
}

func TestSoftwareEngineeringAgent(t *testing.T) {
	t.Run("passes when no code quality issues found", func(t *testing.T) {
		findings := ReviewCode(context.Background(), "diff", stubRunner{out: "[]"})
		require.Len(t, findings, 1)
		assert.Equal(t, Pass, findings[0].Verdict)
	})

	t.Run("returns warn finding for code quality issue", func(t *testing.T) {
		raw := `[{"agent":"software-engineering","verdict":"warn","severity":"low","finding":"naming issue","file":"foo.go:1","fix":"rename it"}]`
		findings := ReviewCode(context.Background(), "diff", stubRunner{out: raw})
		require.Len(t, findings, 1)
		assert.Equal(t, Warn, findings[0].Verdict)
	})

	t.Run("returns block finding for SOLID violation", func(t *testing.T) {
		raw := `[{"agent":"software-engineering","verdict":"block","severity":"high","finding":"SOLID violation","file":"bar.go:5","fix":"extract interface"}]`
		findings := ReviewCode(context.Background(), "diff", stubRunner{out: raw})
		require.Len(t, findings, 1)
		assert.Equal(t, Block, findings[0].Verdict)
	})

	t.Run("skips gracefully on runner error", func(t *testing.T) {
		findings := ReviewCode(context.Background(), "diff", stubRunner{err: errors.New("timeout")})
		require.Len(t, findings, 1)
		assert.Equal(t, Warn, findings[0].Verdict)
		assert.Contains(t, findings[0].Finding, "skipped")
	})

	t.Run("skips gracefully when no runner is configured", func(t *testing.T) {
		findings := ReviewCode(context.Background(), "diff", nil)
		require.Len(t, findings, 1)
		assert.Equal(t, Warn, findings[0].Verdict)
	})

	t.Run("enforces agent name regardless of model response", func(t *testing.T) {
		raw := `[{"agent":"testing-quality","verdict":"warn","severity":"low","finding":"issue","file":"","fix":""}]`
		findings := ReviewCode(context.Background(), "diff", stubRunner{out: raw})
		require.Len(t, findings, 1)
		assert.Equal(t, "software-engineering", findings[0].Agent)
	})

	t.Run("handles markdown-fenced JSON in model response", func(t *testing.T) {
		raw := "```json\n[]\n```"
		findings := ReviewCode(context.Background(), "diff", stubRunner{out: raw})
		require.Len(t, findings, 1)
		assert.Equal(t, Pass, findings[0].Verdict)
	})

	t.Run("skips with warn when model response is unparseable", func(t *testing.T) {
		findings := ReviewCode(context.Background(), "diff", stubRunner{out: "not json"})
		require.Len(t, findings, 1)
		assert.Equal(t, Warn, findings[0].Verdict)
		assert.Contains(t, findings[0].Finding, "unparseable")
	})
}
