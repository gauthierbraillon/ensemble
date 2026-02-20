package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/gauthierbraillon/ensemble/internal/agent"
)

type hookEvent struct {
	ToolName  string    `json:"tool_name"`
	ToolInput toolInput `json:"tool_input"`
}

type toolInput struct {
	FilePath string `json:"file_path"`
}

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Claude Code PreToolUse hook — blocks writes that violate TDD",
	RunE:  runHook,
}

func runHook(_ *cobra.Command, _ []string) error {
	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	var event hookEvent
	if err := json.Unmarshal(raw, &event); err != nil {
		return err
	}
	finding := agent.CheckFileWrite(event.ToolInput.FilePath)
	out, _ := json.Marshal(finding)
	fmt.Println(string(out))
	if finding.Verdict == agent.Block {
		os.Exit(2)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(hookCmd)
}
