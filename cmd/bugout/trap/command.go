package trapcmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bugout-dev/bugout-go/cmd/bugout/cmdutils"
	bugout "github.com/bugout-dev/bugout-go/pkg"
	"github.com/bugout-dev/bugout-go/pkg/spire"
)

func PopulateTrapCommands(cmd *cobra.Command) {
	trapCmd := CreateTrapCommand()
	cmd.AddCommand(trapCmd)
}

func CreateTrapCommand() *cobra.Command {
	var token, journalID, title string
	var tags []string

	trapCmd := &cobra.Command{
		Use:     "trap",
		Short:   "Run commands, capture their output, and add it to a Bugout journal",
		PreRunE: cmdutils.CompositePopulator(cmdutils.TokenArgPopulator, cmdutils.JournalIDArgPopulator),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := RunWrappedCommand(args)
			if err != nil {
				return err
			}

			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			if title == "" {
				title = fmt.Sprintf("Trap: %s", args[0])
			}

			quotedArgs := make([]string, len(args))
			for i, arg := range args {
				quotedArgs[i] = strconv.Quote(arg)
			}

			envvars := os.Environ()
			sort.Strings(envvars)
			quotedEnvvars := make([]string, len(envvars))
			for i, envvar := range envvars {
				quotedEnvvars[i] = fmt.Sprintf("`%s`", envvar)
			}

			var entryContent string = strings.Join([]string{
				fmt.Sprintf("## invocation\n```\n%s\n```\n", strings.Join(quotedArgs, " ")),
				fmt.Sprintf("## code\n`%d`\n", result.ExitCode),
				fmt.Sprintf("## stdout\n```\n%s\n```\n", result.OutBuffer.String()),
				fmt.Sprintf("## stderr\n```\n%s\n```\n", result.ErrBuffer.String()),
				fmt.Sprintf("## env\n- %s\n", strings.Join(quotedEnvvars, "\n- ")),
			}, "\n")

			entry, err := client.Spire.CreateEntry(token, journalID, title, entryContent, tags, spire.EntryContext{})
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(entry)
			return encodeErr
		},
	}

	trapCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	trapCmd.Flags().StringVarP(&journalID, "journal", "j", "", "ID of journal")
	trapCmd.Flags().StringVar(&title, "title", "", "Title of new entry")
	trapCmd.Flags().StringSliceVar(&tags, "tags", []string{}, "Tags to apply to the new entry (as a comma-separated list of strings)")

	return trapCmd
}
