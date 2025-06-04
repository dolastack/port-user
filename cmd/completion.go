package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:
  $ source <(port-user completion bash)

Zsh:
  $ source <(port-user completion zsh)

Fish:
  $ port-user completion fish | source

PowerShell:
  PS> port-user completion powershell | Out-String | Invoke-Expression

To configure your shell to load completions automatically on login:
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Please specify a shell (bash, zsh, fish, powershell)")
			os.Exit(1)
		}
		shell := args[0]

		var err error
		switch shell {
		case "bash":
			err = rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			err = rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			err = rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			err = rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			fmt.Fprintf(os.Stderr, "Unsupported shell: %s\n", shell)
			os.Exit(1)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating completion: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
