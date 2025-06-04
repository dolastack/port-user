package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/dolastack/port-user/internal/core"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "port-user [ports]",
	Short: "A utility to find processes using specific ports",
	Long: `port-user is a cross-platform CLI tool to find which process is using one or more TCP/UDP ports.
Supports Linux, macOS, and Windows.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		portStr := args[0]
		portStrs := strings.Split(portStr, ",")
		var ports []int

		for _, s := range portStrs {
			s = strings.TrimSpace(s)
			port, err := strconv.Atoi(s)
			if err != nil || port < 1 || port > 65535 {
				fmt.Fprintf(os.Stderr, "Invalid port number: %s\n", s)
				os.Exit(1)
			}
			ports = append(ports, port)
		}

		switch runtime.GOOS {
		case "linux", "darwin":
			core.FindPortUnix(ports, false, false)
		case "windows":
			core.FindPortWindows(ports, false, false)
		default:
			fmt.Fprintf(os.Stderr, "Unsupported OS: %s\n", runtime.GOOS)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// GetRootCmd returns the root command for use in other packages
func GetRootCmd() *cobra.Command {
	return rootCmd
}
