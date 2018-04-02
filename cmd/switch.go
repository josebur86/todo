package cmd

import (
    "github.com/spf13/cobra"
)

func init() {
    rootCmd.AddCommand(switchCmd)
}

var switchCmd = &cobra.Command{
    Use: "switch",
    Short: "Switch to another board.",
    Run: runSwitch,
    Args: cobra.MinimumNArgs(1),
}

func runSwitch(cmd *cobra.Command, args []string) {
}
