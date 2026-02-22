package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/gauthierbraillon/ensemble/internal/initcmd"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "Set up ensemble hook in .claude/settings.json",
	Long:    `Creates or updates .claude/settings.json to add the ensemble PreToolUse hook.`,
	Example: `  ensemble init`,
	RunE:    runInit,
}

func runInit(_ *cobra.Command, _ []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	changed, err := initcmd.WriteSettings(dir)
	if err != nil {
		return err
	}
	if changed {
		fmt.Println("Initialised .claude/settings.json — ensemble hook active.")
		fmt.Println("Run: git diff HEAD~1 | ensemble cycle")
	} else {
		fmt.Println("Already configured — nothing changed.")
	}
	if !initcmd.EnsembleOnPath() {
		fmt.Fprintln(os.Stderr, "WARNING: ensemble not found on PATH — hook will not fire until it is installed.")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
