package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dolastack/port-user/internal/core"
)

var (
	jsonOutput bool
	useColor   bool
)

// ANSI color codes
var (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

var portUserCmd = &cobra.Command{
	Use:   "port-user [ports]",
	Short: "Find processes using specified ports",
	Long:  `A utility to find processes using one or more TCP/UDP ports on Unix-like OS and Windows.`,
	Args:  cobra.ExactArgs(1),
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
			core.FindPortUnix(ports, jsonOutput, useColor)
		case "windows":
			core.FindPortWindows(ports, jsonOutput, useColor)
		default:
			fmt.Fprintf(os.Stderr, "Unsupported OS: %s\n", runtime.GOOS)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(portUserCmd)

	portUserCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output result in JSON format")
	portUserCmd.Flags().BoolVarP(&useColor, "color", "c", false, "Enable colorized output")
}
