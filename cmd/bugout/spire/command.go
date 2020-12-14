package spirecmd

import (
	"fmt"

	bugout "github.com/bugout-dev/bugout-go/pkg"
	"github.com/spf13/cobra"
)

func PopulateSpireCommands(cmd *cobra.Command) {
	journalsCmd := CreateJournalsCommand()
	pingCmd := CreatePingCommand()

	cmd.AddCommand(journalsCmd, pingCmd)
}

func CreatePingCommand() *cobra.Command {
	pingCmd := &cobra.Command{
		Use:   "ping-spire",
		Short: "Ping Spire to see if it is active",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			result, err := client.Spire.Ping()
			if err != nil {
				return err
			}

			fmt.Println(result)

			return nil
		},
	}

	return pingCmd
}
