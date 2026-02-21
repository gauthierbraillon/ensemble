package acceptance

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCycleCommandShowsUsageExampleInHelp(t *testing.T) {
	out, err := exec.Command(ensembleBin(t), "cycle", "--help").CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(out), "git diff HEAD~1 | ensemble cycle")
}

func TestHookHelpExplainsItIsCalledByClaudeCode(t *testing.T) {
	out, err := exec.Command(ensembleBin(t), "hook", "--help").CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(out), "settings.json")
}

func TestEnsembleDoesNotExposeInternalShellCompletionCommand(t *testing.T) {
	out, err := exec.Command(ensembleBin(t), "--help").CombinedOutput()
	require.NoError(t, err)
	assert.NotContains(t, string(out), "completion")
}
