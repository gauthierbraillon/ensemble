package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReviewUXSkipsWhenNoExportedAPIChange(t *testing.T) {
	diff := `diff --git a/internal/foo/foo.go b/internal/foo/foo.go
+++ b/internal/foo/foo.go
+func add(a, b int) int { return a + b }`
	findings := ReviewUX(context.Background(), diff, stubRunner{out: "[]"})
	assert.Empty(t, findings)
}

func TestReviewUXPassesOnEmptyArrayWhenExportedChangePresent(t *testing.T) {
	findings := ReviewUX(context.Background(), diffWithExportedFunc(), stubRunner{out: "[]"})
	require.Len(t, findings, 1)
	assert.Equal(t, Pass, findings[0].Verdict)
}

func TestReviewUXReturnsWarnFinding(t *testing.T) {
	raw := `[{"agent":"ux-design","verdict":"warn","severity":"low","finding":"unclear name","file":"api.go:1","fix":"rename to GetUser"}]`
	findings := ReviewUX(context.Background(), diffWithExportedFunc(), stubRunner{out: raw})
	require.Len(t, findings, 1)
	assert.Equal(t, Warn, findings[0].Verdict)
}

func TestReviewUXReturnsBlockFinding(t *testing.T) {
	raw := `[{"agent":"ux-design","verdict":"block","severity":"high","finding":"API breaks convention","file":"api.go:5","fix":"follow REST naming"}]`
	findings := ReviewUX(context.Background(), diffWithExportedFunc(), stubRunner{out: raw})
	require.Len(t, findings, 1)
	assert.Equal(t, Block, findings[0].Verdict)
}

func TestReviewUXSkipsOnNilRunner(t *testing.T) {
	findings := ReviewUX(context.Background(), diffWithExportedFunc(), nil)
	require.Len(t, findings, 1)
	assert.Equal(t, Warn, findings[0].Verdict)
	assert.Contains(t, findings[0].Finding, "skipped")
}

func TestReviewUXSkipsOnRunnerError(t *testing.T) {
	findings := ReviewUX(context.Background(), diffWithExportedFunc(), stubRunner{err: errors.New("timeout")})
	require.Len(t, findings, 1)
	assert.Equal(t, Warn, findings[0].Verdict)
	assert.Contains(t, findings[0].Finding, "skipped")
}

func TestReviewUXEnforcesAgentName(t *testing.T) {
	raw := `[{"agent":"software-engineering","verdict":"warn","severity":"low","finding":"issue","file":"","fix":""}]`
	findings := ReviewUX(context.Background(), diffWithExportedFunc(), stubRunner{out: raw})
	require.Len(t, findings, 1)
	assert.Equal(t, "ux-design", findings[0].Agent)
}

func TestReviewUXHandlesMarkdownFences(t *testing.T) {
	raw := "```json\n[]\n```"
	findings := ReviewUX(context.Background(), diffWithExportedFunc(), stubRunner{out: raw})
	require.Len(t, findings, 1)
	assert.Equal(t, Pass, findings[0].Verdict)
}

func TestHasExportedAPIChangeReturnsTrueForExportedFunc(t *testing.T) {
	diff := "+func Add(a, b int) int { return a + b }"
	assert.True(t, hasExportedAPIChange(diff))
}

func TestHasExportedAPIChangeReturnsFalseForUnexportedFunc(t *testing.T) {
	diff := "+func add(a, b int) int { return a + b }"
	assert.False(t, hasExportedAPIChange(diff))
}

func TestHasExportedAPIChangeReturnsFalseForTestFiles(t *testing.T) {
	diff := `+++ b/foo_test.go
+func TestAdd(t *testing.T) {}
+func Add(a, b int) int { return a + b }`
	assert.False(t, hasExportedAPIChange(diff))
}

func TestHasExportedAPIChangeReturnsTrueForExportedType(t *testing.T) {
	diff := "+type Config struct { Host string }"
	assert.True(t, hasExportedAPIChange(diff))
}

func TestHasExportedAPIChangeReturnsTrueForExportedConst(t *testing.T) {
	diff := "+const MaxRetries = 3"
	assert.True(t, hasExportedAPIChange(diff))
}

func TestHasExportedAPIChangeReturnsFalseForRemovedOnlyLines(t *testing.T) {
	diff := "-func Remove() {}"
	assert.False(t, hasExportedAPIChange(diff))
}

func diffWithExportedFunc() string {
	return `diff --git a/internal/foo/foo.go b/internal/foo/foo.go
+++ b/internal/foo/foo.go
+func Add(a, b int) int { return a + b }`
}
