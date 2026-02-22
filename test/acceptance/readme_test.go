package acceptance

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocumentation(t *testing.T) {
	t.Run("README exists", func(t *testing.T) {
		_, err := os.Stat("../../README.md")
		require.NoError(t, err)
	})

	t.Run("README documents install, hook setup, cycle usage, and minimum CD reference", func(t *testing.T) {
		content, err := os.ReadFile("../../README.md")
		require.NoError(t, err)
		s := string(content)
		assert.Contains(t, s, "ensemble")
		assert.Contains(t, s, "install")
		assert.Contains(t, s, "hook")
		assert.Contains(t, s, "minimumcd.org")
	})
}
