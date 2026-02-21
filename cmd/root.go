package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ensemble",
	Short: "AI agent team — quality-enforced TDD mob programming",
	Long: `ensemble runs an AI agent team that enforces XP/CD discipline on every code change.

You direct. Agents enforce. Quality gates are non-negotiable.`,
	RunE: runSession,
}

func runSession(_ *cobra.Command, _ []string) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("ensemble> ")
	for scanner.Scan() {
		fmt.Print("ensemble> ")
	}
	return nil
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
