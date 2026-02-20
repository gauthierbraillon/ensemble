package acceptance

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeployScriptExists(t *testing.T) {
	info, err := os.Stat("../../deploy.sh")
	require.NoError(t, err)
	assert.NotZero(t, info.Mode()&0111, "deploy.sh must be executable")
}

func TestDeployScriptRunsGatesBeforePush(t *testing.T) {
	content, err := os.ReadFile("../../deploy.sh")
	require.NoError(t, err)
	s := string(content)
	assert.Contains(t, s, "make ci")
	assert.Contains(t, s, "git push")
	assert.Contains(t, s, "semantic-release")
}
