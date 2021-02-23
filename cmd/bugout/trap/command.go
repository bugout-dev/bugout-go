package trapcmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bugout-dev/bugout-go/cmd/bugout/cmdutils"
	bugout "github.com/bugout-dev/bugout-go/pkg"
)

func PopulateTrapCommands(cmd *cobra.Command) {
	trapCmd := CreateTrapCommand()
	cmd.AddCommand(trapCmd)
}

func CreateTrapCommand() *cobra.Command {
	var token, journalID, title string
	var tags []string
	var showEnv bool

	trapCmd := &cobra.Command{
		Use:   "trap",
		Short: "Run commands, capture their output, and add it to a Bugout journal",
		Long: `Wraps a command, waits for it to complete, and then adds the result to a Bugout journal.

Specify the wrapped command using "--" followed by the command:
	bugout trap [flags] -- <command>
`,
		Args: func(cmd *cobra.Command, args []string) error {
			populateMissingArgsFromEnv := cmdutils.CompositePopulator(cmdutils.TokenArgPopulator, cmdutils.JournalIDArgPopulator)
			err := populateMissingArgsFromEnv(cmd, args)
			if err != nil {
				return err
			}
			if len(args) == 0 {
				return errors.New("You must pass this command an invocation")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := RunWrappedCommand(cmd, args)
			if err != nil {
				return err
			}

			entry := Render(args, result, title, tags, showEnv)

			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			response, err := client.Spire.CreateEntry(token, journalID, entry.Title, entry.Content, entry.Tags, entry.Context)
			if err != nil {
				return err
			}

			cmd.ErrOrStderr().Write([]byte(fmt.Sprintf("\n\nBugout entry created at: %s/journals/%s/%s\n", cmdutils.BugoutURL(), journalID, response.Id)))

			if result.ExitCode > 0 {
				os.Exit(result.ExitCode)
			}
			return nil
		},
	}

	trapCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	trapCmd.Flags().StringVarP(&journalID, "journal", "j", "", "ID of journal")
	trapCmd.Flags().StringVarP(&title, "title", "T", "", "Title of new entry")
	trapCmd.Flags().StringSliceVar(&tags, "tags", []string{}, "Tags to apply to the new entry (as a comma-separated list of strings)")
	trapCmd.Flags().BoolVarP(&showEnv, "env", "e", false, "Set this flag to dump the values of your current environment variables")

	return trapCmd
}
