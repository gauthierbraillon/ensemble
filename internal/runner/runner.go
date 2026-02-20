package runner

import (
	"context"
	"time"
)

type Runner interface {
	Run(ctx context.Context, prompt string) (string, error)
}

type Config struct {
	Binary   string
	Model    string
	Timeout  time.Duration
	MaxBytes int64
	ExtraEnv []string
}

var ValidModels = map[string]bool{
	"claude-opus-4-6":           true,
	"claude-sonnet-4-6":         true,
	"claude-haiku-4-5-20251001": true,
}
