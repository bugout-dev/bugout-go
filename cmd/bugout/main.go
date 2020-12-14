package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	broodcmd "github.com/bugout-dev/bugout-go/cmd/bugout/brood"
	spirecmd "github.com/bugout-dev/bugout-go/cmd/bugout/spire"
	bugout "github.com/bugout-dev/bugout-go/pkg"
)

func CreateBugoutCommand() *cobra.Command {
	bugoutCmd := &cobra.Command{
		Use:   "bugout",
		Short: "Interact with Bugout from your command line",
		Long: `Bugout: The knowledge base for software teams

The bugout utility lets you interact with your Bugout resources from your command line.`,
		Version: bugout.Version,
	}

	broodcmd.PopulateBroodCommands(bugoutCmd)
	spirecmd.PopulateSpireCommands(bugoutCmd)

	completionCmd := CreateBugoutCompletionCommand()
	bugoutCmd.AddCommand(completionCmd)

	return bugoutCmd
}

func CreateBugoutCompletionCommand() *cobra.Command {
	completionCmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate shell completion scripts for the bugout tool",
		Long: `Generate shell completion scripts for the bugout tool

You can source these generated scripts or add them in the appropriate completion directories to
unlock bugout command completion for your shell.

For example, to activate bash completions:
	$ . <(bugout completion bash)`,
	}

	bashCompletionCmd := &cobra.Command{
		Use:   "bash",
		Short: "bash completions for bugout",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Root().GenBashCompletion(cmd.OutOrStdout())
		},
	}

	zshCompletionCmd := &cobra.Command{
		Use:   "zsh",
		Short: "zsh completions for bugout",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Root().GenZshCompletion(cmd.OutOrStdout())
		},
	}

	fishCompletionCmd := &cobra.Command{
		Use:   "fish",
		Short: "fish completions for bugout",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
		},
	}

	powershellCompletionCmd := &cobra.Command{
		Use:   "powershell",
		Short: "powershell completions for bugout",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Root().GenPowerShellCompletion(cmd.OutOrStdout())
		},
	}

	completionCmd.AddCommand(bashCompletionCmd, zshCompletionCmd, fishCompletionCmd, powershellCompletionCmd)

	return completionCmd
}

func main() {
	bugCmd := CreateBugoutCommand()
	err := bugCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
