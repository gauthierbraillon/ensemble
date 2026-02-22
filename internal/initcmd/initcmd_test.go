package initcmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gauthierbraillon/ensemble/internal/initcmd"
)

func TestWriteSettings(t *testing.T) {
	t.Run("creates .claude/settings.json when absent", func(t *testing.T) {
		dir := t.TempDir()
		changed, err := initcmd.WriteSettings(dir)
		require.NoError(t, err)
		assert.True(t, changed)
		data, err := os.ReadFile(filepath.Join(dir, ".claude", "settings.json"))
		require.NoError(t, err)
		assert.Contains(t, string(data), "ensemble hook")
	})

	t.Run("is idempotent — returns unchanged on second call", func(t *testing.T) {
		dir := t.TempDir()
		_, _ = initcmd.WriteSettings(dir)
		changed, err := initcmd.WriteSettings(dir)
		require.NoError(t, err)
		assert.False(t, changed)
	})

	t.Run("merges hook when settings.json exists with other keys", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(dir, ".claude"), 0755))
		require.NoError(t, os.WriteFile(
			filepath.Join(dir, ".claude", "settings.json"),
			[]byte(`{"someOtherKey": true}`),
			0644,
		))
		changed, err := initcmd.WriteSettings(dir)
		require.NoError(t, err)
		assert.True(t, changed)
		data, err := os.ReadFile(filepath.Join(dir, ".claude", "settings.json"))
		require.NoError(t, err)
		assert.Contains(t, string(data), "someOtherKey")
		assert.Contains(t, string(data), "ensemble hook")
	})

	t.Run("reports already configured when hook is present", func(t *testing.T) {
		dir := t.TempDir()
		_, err := initcmd.WriteSettings(dir)
		require.NoError(t, err)
		changed, err := initcmd.WriteSettings(dir)
		require.NoError(t, err)
		assert.False(t, changed)
	})
}
