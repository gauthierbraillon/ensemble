package acceptance

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestREADMEExists(t *testing.T) {
	_, err := os.Stat("../../README.md")
	require.NoError(t, err)
}

func TestREADMEContainsEssentialSections(t *testing.T) {
	content, err := os.ReadFile("../../README.md")
	require.NoError(t, err)
	s := string(content)

	assert.Contains(t, s, "ensemble")
	assert.Contains(t, s, "install")
	assert.Contains(t, s, "hook")
	assert.Contains(t, s, "minimumcd.org")
}
