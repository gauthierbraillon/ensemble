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

func TestReviewCodePassesOnEmptyArray(t *testing.T) {
	findings := ReviewCode(context.Background(), "diff", stubRunner{out: "[]"})
	require.Len(t, findings, 1)
	assert.Equal(t, Pass, findings[0].Verdict)
}

func TestReviewCodeReturnsWarnFinding(t *testing.T) {
	raw := `[{"agent":"software-engineering","verdict":"warn","severity":"low","finding":"naming issue","file":"foo.go:1","fix":"rename it"}]`
	findings := ReviewCode(context.Background(), "diff", stubRunner{out: raw})
	require.Len(t, findings, 1)
	assert.Equal(t, Warn, findings[0].Verdict)
}

func TestReviewCodeReturnsBlockFinding(t *testing.T) {
	raw := `[{"agent":"software-engineering","verdict":"block","severity":"high","finding":"SOLID violation","file":"bar.go:5","fix":"extract interface"}]`
	findings := ReviewCode(context.Background(), "diff", stubRunner{out: raw})
	require.Len(t, findings, 1)
	assert.Equal(t, Block, findings[0].Verdict)
}

func TestReviewCodeSkipsOnRunnerError(t *testing.T) {
	findings := ReviewCode(context.Background(), "diff", stubRunner{err: errors.New("timeout")})
	require.Len(t, findings, 1)
	assert.Equal(t, Warn, findings[0].Verdict)
	assert.Contains(t, findings[0].Finding, "skipped")
}

func TestReviewCodeSkipsOnNilRunner(t *testing.T) {
	findings := ReviewCode(context.Background(), "diff", nil)
	require.Len(t, findings, 1)
	assert.Equal(t, Warn, findings[0].Verdict)
}

func TestReviewCodeEnforcesAgentName(t *testing.T) {
	raw := `[{"agent":"testing-quality","verdict":"warn","severity":"low","finding":"issue","file":"","fix":""}]`
	findings := ReviewCode(context.Background(), "diff", stubRunner{out: raw})
	require.Len(t, findings, 1)
	assert.Equal(t, "software-engineering", findings[0].Agent)
}

func TestReviewCodeHandlesMarkdownFences(t *testing.T) {
	raw := "```json\n[]\n```"
	findings := ReviewCode(context.Background(), "diff", stubRunner{out: raw})
	require.Len(t, findings, 1)
	assert.Equal(t, Pass, findings[0].Verdict)
}

func TestReviewCodeSkipsOnUnparseableResponse(t *testing.T) {
	findings := ReviewCode(context.Background(), "diff", stubRunner{out: "not json"})
	require.Len(t, findings, 1)
	assert.Equal(t, Warn, findings[0].Verdict)
	assert.Contains(t, findings[0].Finding, "unparseable")
}
