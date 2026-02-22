package acceptance

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInteractiveSession(t *testing.T) {
	t.Run("opens interactive prompt when run with no arguments", func(t *testing.T) {
		cmd := exec.Command(ensembleBin(t))
		cmd.Stdin = strings.NewReader("")
		out, err := cmd.CombinedOutput()
		require.NoError(t, err)
		assert.Contains(t, string(out), "ensemble> ")
		assert.NotContains(t, string(out), "Available Commands:")
	})
}
