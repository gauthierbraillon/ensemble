package acceptance

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLIHelp(t *testing.T) {
	t.Run("cycle --help shows usage example", func(t *testing.T) {
		out, err := exec.Command(ensembleBin(t), "cycle", "--help").CombinedOutput()
		require.NoError(t, err)
		assert.Contains(t, string(out), "git diff HEAD~1 | ensemble cycle")
	})

	t.Run("hook --help explains it is called by Claude Code via settings.json", func(t *testing.T) {
		out, err := exec.Command(ensembleBin(t), "hook", "--help").CombinedOutput()
		require.NoError(t, err)
		assert.Contains(t, string(out), "settings.json")
	})

	t.Run("does not expose internal shell completion command", func(t *testing.T) {
		out, err := exec.Command(ensembleBin(t), "--help").CombinedOutput()
		require.NoError(t, err)
		assert.NotContains(t, string(out), "completion")
	})
}
