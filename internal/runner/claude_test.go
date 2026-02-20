package runner_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gauthierbraillon/ensemble/internal/runner"
)

func TestMain(m *testing.M) {
	switch os.Getenv("RUNNER_SUBPROC") {
	case "echo_stdin":
		_, _ = io.Copy(os.Stdout, os.Stdin)
		os.Exit(0)
	case "sleep":
		time.Sleep(time.Hour)
		os.Exit(0)
	case "big_output":
		fmt.Print(strings.Repeat("x", 200))
		os.Exit(0)
	case "print_env":
		for _, kv := range os.Environ() {
			fmt.Println(kv)
		}
		os.Exit(0)
	case "exit_fail":
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func selfBin(t *testing.T) string {
	t.Helper()
	bin, err := os.Executable()
	require.NoError(t, err)
	return bin
}

func cfg(t *testing.T, subproc string, maxBytes int64) runner.Config {
	t.Helper()
	return runner.Config{
		Binary:   selfBin(t),
		Model:    "",
		Timeout:  2 * time.Second,
		MaxBytes: maxBytes,
		ExtraEnv: []string{"RUNNER_SUBPROC=" + subproc},
	}
}

func TestRunnerPassesPromptViaStdin(t *testing.T) {
	r, err := runner.New(cfg(t, "echo_stdin", 1024))
	require.NoError(t, err)
	out, err := r.Run(context.Background(), "hello")
	require.NoError(t, err)
	assert.Equal(t, "hello", out)
}

func TestRunnerEnforcesTimeout(t *testing.T) {
	c := cfg(t, "sleep", 1024)
	c.Timeout = 100 * time.Millisecond
	r, err := runner.New(c)
	require.NoError(t, err)
	_, err = r.Run(context.Background(), "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
}

func TestRunnerBlocksOutputExceedingMaxBytes(t *testing.T) {
	r, err := runner.New(cfg(t, "big_output", 100))
	require.NoError(t, err)
	_, err = r.Run(context.Background(), "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "output exceeded")
}

func TestRunnerSanitizesEnvironment(t *testing.T) {
	t.Setenv("AWS_SECRET_ACCESS_KEY", "supersecret")
	t.Setenv("GITHUB_TOKEN", "ghp_fake")
	r, err := runner.New(cfg(t, "print_env", 4096))
	require.NoError(t, err)
	out, err := r.Run(context.Background(), "")
	require.NoError(t, err)
	assert.NotContains(t, out, "supersecret")
	assert.NotContains(t, out, "ghp_fake")
}

func TestRunnerReturnsErrorOnNonZeroExit(t *testing.T) {
	r, err := runner.New(cfg(t, "exit_fail", 1024))
	require.NoError(t, err)
	_, err = r.Run(context.Background(), "")
	require.Error(t, err)
}

func TestNewRejectsUnknownModel(t *testing.T) {
	c := runner.Config{Binary: "claude", Model: "gpt-4", Timeout: time.Second, MaxBytes: 1024}
	_, err := runner.New(c)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid model")
}

func TestNewAcceptsEmptyModelForTesting(t *testing.T) {
	c := runner.Config{Binary: selfBin(t), Model: "", Timeout: time.Second, MaxBytes: 1024}
	_, err := runner.New(c)
	require.NoError(t, err)
}
