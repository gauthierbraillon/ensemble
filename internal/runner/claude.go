package runner

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var allowedEnvKeys = []string{"HOME", "PATH", "TERM", "LANG"}
var allowedEnvPrefixes = []string{"LC_"}

type ClaudeRunner struct {
	cfg        Config
	binaryPath string
}

func New(cfg Config) (*ClaudeRunner, error) {
	if cfg.Model != "" && !ValidModels[cfg.Model] {
		return nil, fmt.Errorf("invalid model: %q", cfg.Model)
	}
	resolved, err := exec.LookPath(cfg.Binary)
	if err != nil {
		return nil, fmt.Errorf("binary %q not found: %w", cfg.Binary, err)
	}
	return &ClaudeRunner{cfg: cfg, binaryPath: resolved}, nil
}

func (r *ClaudeRunner) Run(ctx context.Context, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.Timeout)
	defer cancel()

	args := []string{"--print"}
	if r.cfg.Model != "" {
		args = append(args, "--model", r.cfg.Model)
	}

	cmd := exec.CommandContext(ctx, r.binaryPath, args...) // #nosec G204
	cmd.Stdin = strings.NewReader(prompt)
	cmd.Env = append(safeEnv(os.Environ()), r.cfg.ExtraEnv...)

	lb := &limitedBuffer{max: r.cfg.MaxBytes}
	cmd.Stdout = lb

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("runner: timeout after %s", r.cfg.Timeout)
		}
		if lb.exceeded {
			return "", fmt.Errorf("runner: output exceeded %d bytes", r.cfg.MaxBytes)
		}
		return "", fmt.Errorf("runner: %w", err)
	}

	return lb.buf.String(), nil
}

func safeEnv(env []string) []string {
	var safe []string
	for _, kv := range env {
		if isAllowed(kv) {
			safe = append(safe, kv)
		}
	}
	return safe
}

func isAllowed(kv string) bool {
	for _, key := range allowedEnvKeys {
		if strings.HasPrefix(kv, key+"=") {
			return true
		}
	}
	for _, prefix := range allowedEnvPrefixes {
		if strings.HasPrefix(kv, prefix) {
			return true
		}
	}
	return false
}

type limitedBuffer struct {
	buf      bytes.Buffer
	max      int64
	written  int64
	exceeded bool
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	if b.written+int64(len(p)) > b.max {
		b.exceeded = true
		return 0, fmt.Errorf("output cap exceeded")
	}
	n, err := b.buf.Write(p)
	b.written += int64(n)
	return n, err
}
