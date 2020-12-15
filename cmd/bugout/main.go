package main

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			viper.SetConfigName("bugout")
			viper.SetConfigType("toml")
			viper.AddConfigPath("./")
			homeDir, homeDirErr := os.UserHomeDir()
			if homeDirErr == nil {
				viper.AddConfigPath(path.Join(homeDir, ".bugout"))
				viper.AddConfigPath(homeDir)
			}
			return viper.ReadInConfig()
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return viper.WriteConfig()
		},
	}

	broodcmd.PopulateBroodCommands(bugoutCmd)
	spirecmd.PopulateSpireCommands(bugoutCmd)

	completionCmd := CreateBugoutCompletionCommand()
	stateCmd := CreateBugoutStateCommand()
	bugoutCmd.AddCommand(completionCmd, stateCmd)

	return bugoutCmd
}

func CreateBugoutStateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "Maintain bugout state across invocations",
		Long: `Operations that allow you to maintain and inspect bugout state between bugout invocations.

This is used to store things like bugout access tokens and active journal IDs.`,
	}

	initCmd := CreateBugoutStateInitCommand()
	currentCmd := CreateBugoutStateCurrentCommand()
	setCmd := CreateBugoutStateSetCommand()
	cmd.AddCommand(initCmd, currentCmd, setCmd)

	return cmd
}

func CreateBugoutStateInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a bugout.toml configuration file in $HOME/.bugout/bugout.toml",
		RunE: func(cmd *cobra.Command, args []string) error {
			userHome, userHomeErr := os.UserHomeDir()
			if userHomeErr != nil {
				return userHomeErr
			}

			bugoutDirPath := path.Join(userHome, ".bugout")
			stat, statErr := os.Stat(bugoutDirPath)
			nonexistence := os.IsNotExist(statErr)
			if !nonexistence && !stat.IsDir() {
				return fmt.Errorf("%s exists but is not a directory", bugoutDirPath)
			}
			if nonexistence {
				os.Mkdir(bugoutDirPath, 0755)
			}

			configFilePath := path.Join(bugoutDirPath, "bugout.toml")
			configFile, configFileErr := os.OpenFile(configFilePath, os.O_RDONLY|os.O_CREATE, 0644)
			if configFileErr != nil {
				return configFileErr
			}
			return configFile.Close()
		},
	}

	return cmd
}

func CreateBugoutStateCurrentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Show the current bugout state",
		Run: func(cmd *cobra.Command, args []string) {
			for _, key := range viper.AllKeys() {
				value := viper.GetString(key)
				fmt.Printf("%s: %s\n", key, value)
			}
		},
	}

	return cmd
}

func CreateBugoutStateSetCommand() *cobra.Command {
	var key, value string
	cmd := &cobra.Command{
		Use:   "set KEY VALUE",
		Short: "Set a key-value pair in the bugout state",
		Long:  "Set a key-value pair in the bugout state\nValid keys: access_token,user_id,journal_id",
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			key = args[0]
			value = args[1]
			if key != "access_token" && key != "user_id" && key != "journal_id" {
				return fmt.Errorf("Invalid key: %s. Valid keys: access_token,user_id,journal_id", key)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			viper.Set(key, value)
		},
	}

	return cmd
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
