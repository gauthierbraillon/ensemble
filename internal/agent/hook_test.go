package agent_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gauthierbraillon/ensemble/internal/agent"
)

func TestCheckFileWritePassesForNonGoFile(t *testing.T) {
	f := agent.CheckFileWrite("/project/README.md")
	assert.Equal(t, agent.Pass, f.Verdict)
}

func TestCheckFileWritePassesForTestFile(t *testing.T) {
	f := agent.CheckFileWrite("/project/internal/foo/foo_test.go")
	assert.Equal(t, agent.Pass, f.Verdict)
}

func TestCheckFileWritePassesWhenTestFileExistsOnDisk(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "foo_test.go")
	require.NoError(t, os.WriteFile(testFile, []byte("package foo_test"), 0600))

	f := agent.CheckFileWrite(filepath.Join(dir, "foo.go"))
	assert.Equal(t, agent.Pass, f.Verdict)
}

func TestCheckFileWriteBlocksWhenNoTestFileOnDisk(t *testing.T) {
	dir := t.TempDir()
	f := agent.CheckFileWrite(filepath.Join(dir, "foo.go"))
	assert.Equal(t, agent.Block, f.Verdict)
	assert.Equal(t, agent.Critical, f.Severity)
}
