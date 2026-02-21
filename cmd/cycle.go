package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/gauthierbraillon/ensemble/internal/agent"
	"github.com/gauthierbraillon/ensemble/internal/runner"
)

var cycleCmd = &cobra.Command{
	Use:   "cycle",
	Short: "Enforce RED→GREEN→REFACTOR→DEPLOY on a diff",
	Long: `Reads a unified diff from stdin and runs TDD and software engineering agents against it.

Each finding prints as one JSON line. Exits 1 if any verdict is "block".`,
	Example: `  git diff HEAD~1 | ensemble cycle
  git diff HEAD   | ensemble cycle
  ensemble cycle  < my.patch`,
	RunE: runCycle,
}

func runCycle(_ *cobra.Command, _ []string) error {
	diff, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	var findings []agent.Finding
	findings = append(findings, agent.ReviewDiff(string(diff))...)
	findings = append(findings, agent.ReviewCode(context.Background(), string(diff), sweRunner())...)
	blocked := false
	for _, f := range findings {
		line, _ := json.Marshal(f)
		fmt.Println(string(line))
		if f.Verdict == agent.Block {
			blocked = true
		}
	}
	if blocked {
		os.Exit(1)
	}
	return nil
}

func sweRunner() runner.Runner {
	if os.Getenv("ANTHROPIC_API_KEY") == "" {
		return nil
	}
	r, err := runner.New(runner.Config{
		Binary:   "claude",
		Model:    "claude-haiku-4-5-20251001",
		Timeout:  30 * time.Second,
		MaxBytes: 32 * 1024,
	})
	if err != nil {
		return nil
	}
	return r
}

func init() {
	rootCmd.AddCommand(cycleCmd)
}
