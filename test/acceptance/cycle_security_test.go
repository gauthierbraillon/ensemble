package acceptance

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCycleIncludesSecurityAgentInOutput(t *testing.T) {
	cmd := exec.Command(ensembleBin(t), "cycle")
	cmd.Stdin = strings.NewReader(diffWithTest())
	cmd.Env = envWithout(os.Environ(), "ANTHROPIC_API_KEY")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "expected exit 0: %s", out)

	findings := parseFindings(t, out)
	hasSecurityAgent := false
	for _, f := range findings {
		if f["agent"] == "security" {
			hasSecurityAgent = true
		}
	}
	assert.True(t, hasSecurityAgent, "no security agent finding in output")
}

func TestCycleSecurityAgentRunsOfflineWithWarnFallback(t *testing.T) {
	cmd := exec.Command(ensembleBin(t), "cycle")
	cmd.Stdin = strings.NewReader(diffWithTest())
	cmd.Env = envWithout(os.Environ(), "ANTHROPIC_API_KEY")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "expected exit 0: %s", out)

	findings := parseFindings(t, out)
	for _, f := range findings {
		if f["agent"] == "security" {
			assert.Equal(t, "warn", f["verdict"])
		}
	}
}
