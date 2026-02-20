package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/gauthierbraillon/ensemble/internal/agent"
)

var cycleCmd = &cobra.Command{
	Use:   "cycle",
	Short: "Enforce RED→GREEN→REFACTOR→DEPLOY on a diff read from stdin",
	RunE:  runCycle,
}

func runCycle(_ *cobra.Command, _ []string) error {
	diff, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	findings := agent.ReviewDiff(string(diff))
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

func init() {
	rootCmd.AddCommand(cycleCmd)
}
