package acceptance

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ensembleBinAbs(t *testing.T) string {
	t.Helper()
	abs, err := filepath.Abs(ensembleBin(t))
	require.NoError(t, err)
	return abs
}

func TestEnsembleInit(t *testing.T) {
	t.Run("creates .claude/settings.json in the project directory", func(t *testing.T) {
		dir := t.TempDir()
		cmd := exec.Command(ensembleBinAbs(t), "init")
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "unexpected error: %s", out)
		data, err := os.ReadFile(filepath.Join(dir, ".claude", "settings.json"))
		require.NoError(t, err, ".claude/settings.json not created")
		assert.Contains(t, string(data), "ensemble hook")
	})

	t.Run("output confirms ensemble hook is active", func(t *testing.T) {
		dir := t.TempDir()
		cmd := exec.Command(ensembleBinAbs(t), "init")
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "unexpected error: %s", out)
		combined := string(out)
		assert.Contains(t, combined, "Initialised")
		assert.Contains(t, combined, "ensemble hook")
		assert.Contains(t, combined, "git diff HEAD~1 | ensemble cycle")
	})

	t.Run("is idempotent — safe to run twice without changes", func(t *testing.T) {
		dir := t.TempDir()
		bin := ensembleBinAbs(t)
		cmd1 := exec.Command(bin, "init")
		cmd1.Dir = dir
		_, err := cmd1.CombinedOutput()
		require.NoError(t, err)
		data1, err := os.ReadFile(filepath.Join(dir, ".claude", "settings.json"))
		require.NoError(t, err)

		cmd2 := exec.Command(bin, "init")
		cmd2.Dir = dir
		out2, err := cmd2.CombinedOutput()
		require.NoError(t, err, "second run failed: %s", out2)
		assert.Contains(t, string(out2), "Already configured")
		data2, err := os.ReadFile(filepath.Join(dir, ".claude", "settings.json"))
		require.NoError(t, err)
		assert.Equal(t, data1, data2, "file should be unchanged")
	})

	t.Run("merges with existing settings.json preserving other content", func(t *testing.T) {
		dir := t.TempDir()
		claudeDir := filepath.Join(dir, ".claude")
		require.NoError(t, os.MkdirAll(claudeDir, 0755))
		require.NoError(t, os.WriteFile(
			filepath.Join(claudeDir, "settings.json"),
			[]byte(`{"someOtherKey": true}`),
			0644,
		))
		cmd := exec.Command(ensembleBinAbs(t), "init")
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "unexpected error: %s", out)
		data, err := os.ReadFile(filepath.Join(claudeDir, "settings.json"))
		require.NoError(t, err)
		content := string(data)
		assert.Contains(t, content, "someOtherKey")
		assert.Contains(t, content, "ensemble hook")
	})

	t.Run("warns when ensemble is not on PATH", func(t *testing.T) {
		dir := t.TempDir()
		cmd := exec.Command(ensembleBinAbs(t), "init")
		cmd.Dir = dir
		cmd.Env = append(envWithout(os.Environ(), "PATH"), "PATH=/usr/bin:/bin")
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "unexpected error: %s", out)
		combined := string(out)
		assert.Contains(t, combined, "WARNING")
		assert.Contains(t, combined, "not found on PATH")
	})
}
