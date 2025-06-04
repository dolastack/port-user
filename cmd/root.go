package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "port-user",
	Short: "A utility to find processes using specific ports",
	Long: `port-user is a cross-platform CLI tool to find which process is using one or more TCP/UDP ports.
Supports Linux, macOS, and Windows.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}
