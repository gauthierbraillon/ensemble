package acceptance

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func hookEvent(toolName, filePath string) string {
	return fmt.Sprintf(`{"tool_name":%q,"tool_input":{"file_path":%q}}`, toolName, filePath)
}

func TestHookIgnoresNonGoFiles(t *testing.T) {
	cmd := exec.Command(ensembleBin(t), "hook")
	cmd.Stdin = strings.NewReader(hookEvent("Write", "/project/README.md"))
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err, "unexpected block: %s", out)
}

func TestHookAllowsWritingTestFilesDirectly(t *testing.T) {
	cmd := exec.Command(ensembleBin(t), "hook")
	cmd.Stdin = strings.NewReader(hookEvent("Write", "/project/internal/foo/foo_test.go"))
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err, "unexpected block: %s", out)
}

func TestHookAllowsWritingImplementationWhenTestAlreadyExists(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "foo_test.go"), []byte("package foo_test"), 0600))

	cmd := exec.Command(ensembleBin(t), "hook")
	cmd.Stdin = strings.NewReader(hookEvent("Write", filepath.Join(dir, "foo.go")))
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err, "unexpected block: %s", out)
}

func TestHookStopsClaudeCodeFromWritingCodeBeforeATest(t *testing.T) {
	dir := t.TempDir()
	cmd := exec.Command(ensembleBin(t), "hook")
	cmd.Stdin = strings.NewReader(hookEvent("Write", filepath.Join(dir, "bar.go")))
	out, _ := cmd.CombinedOutput()

	assert.Equal(t, 2, cmd.ProcessState.ExitCode(), "hook must exit 2 to block Claude Code: %s", out)

	var f map[string]interface{}
	require.NoError(t, json.Unmarshal(out, &f))
	assert.Equal(t, "block", f["verdict"])
}

func TestHookEnforcesToDDOnEditsNotJustNewFiles(t *testing.T) {
	dir := t.TempDir()
	cmd := exec.Command(ensembleBin(t), "hook")
	cmd.Stdin = strings.NewReader(hookEvent("Edit", filepath.Join(dir, "bar.go")))
	out, _ := cmd.CombinedOutput()

	assert.Equal(t, 2, cmd.ProcessState.ExitCode(), "hook must exit 2 for Edit without test: %s", out)
}
