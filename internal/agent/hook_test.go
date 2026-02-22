package agent_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gauthierbraillon/ensemble/internal/agent"
)

func TestTDDHookEnforcement(t *testing.T) {
	t.Run("passes for non-Go files", func(t *testing.T) {
		f := agent.CheckFileWrite("/project/README.md")
		assert.Equal(t, agent.Pass, f.Verdict)
	})

	t.Run("passes for test files — writing tests is always allowed", func(t *testing.T) {
		f := agent.CheckFileWrite("/project/internal/foo/foo_test.go")
		assert.Equal(t, agent.Pass, f.Verdict)
	})

	t.Run("passes when matching test file already exists on disk", func(t *testing.T) {
		dir := t.TempDir()
		testFile := filepath.Join(dir, "foo_test.go")
		require.NoError(t, os.WriteFile(testFile, []byte("package foo_test"), 0600))

		f := agent.CheckFileWrite(filepath.Join(dir, "foo.go"))
		assert.Equal(t, agent.Pass, f.Verdict)
	})

	t.Run("blocks with critical severity when no test file exists on disk", func(t *testing.T) {
		dir := t.TempDir()
		f := agent.CheckFileWrite(filepath.Join(dir, "foo.go"))
		assert.Equal(t, agent.Block, f.Verdict)
		assert.Equal(t, agent.Critical, f.Severity)
	})
}
