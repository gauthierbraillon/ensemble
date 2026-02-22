package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReviewSecurityPassesOnEmptyArray(t *testing.T) {
	findings := ReviewSecurity(context.Background(), "diff", stubRunner{out: "[]"})
	require.Len(t, findings, 1)
	assert.Equal(t, Pass, findings[0].Verdict)
}

func TestReviewSecurityReturnsWarnFinding(t *testing.T) {
	raw := `[{"agent":"security","verdict":"warn","severity":"medium","finding":"SQL injection risk","file":"db.go:10","fix":"use parameterized query"}]`
	findings := ReviewSecurity(context.Background(), "diff", stubRunner{out: raw})
	require.Len(t, findings, 1)
	assert.Equal(t, Warn, findings[0].Verdict)
}

func TestReviewSecurityReturnsBlockFinding(t *testing.T) {
	raw := `[{"agent":"security","verdict":"block","severity":"critical","finding":"hardcoded secret","file":"config.go:3","fix":"use env var"}]`
	findings := ReviewSecurity(context.Background(), "diff", stubRunner{out: raw})
	require.Len(t, findings, 1)
	assert.Equal(t, Block, findings[0].Verdict)
}

func TestReviewSecuritySkipsOnNilRunner(t *testing.T) {
	findings := ReviewSecurity(context.Background(), "diff", nil)
	require.Len(t, findings, 1)
	assert.Equal(t, Warn, findings[0].Verdict)
	assert.Contains(t, findings[0].Finding, "skipped")
}

func TestReviewSecuritySkipsOnRunnerError(t *testing.T) {
	findings := ReviewSecurity(context.Background(), "diff", stubRunner{err: errors.New("timeout")})
	require.Len(t, findings, 1)
	assert.Equal(t, Warn, findings[0].Verdict)
	assert.Contains(t, findings[0].Finding, "skipped")
}

func TestReviewSecurityEnforcesAgentName(t *testing.T) {
	raw := `[{"agent":"software-engineering","verdict":"warn","severity":"low","finding":"issue","file":"","fix":""}]`
	findings := ReviewSecurity(context.Background(), "diff", stubRunner{out: raw})
	require.Len(t, findings, 1)
	assert.Equal(t, "security", findings[0].Agent)
}

func TestReviewSecurityHandlesMarkdownFences(t *testing.T) {
	raw := "```json\n[]\n```"
	findings := ReviewSecurity(context.Background(), "diff", stubRunner{out: raw})
	require.Len(t, findings, 1)
	assert.Equal(t, Pass, findings[0].Verdict)
}

func TestReviewSecuritySkipsOnUnparseableResponse(t *testing.T) {
	findings := ReviewSecurity(context.Background(), "diff", stubRunner{out: "not json"})
	require.Len(t, findings, 1)
	assert.Equal(t, Warn, findings[0].Verdict)
	assert.Contains(t, findings[0].Finding, "unparseable")
}
