package spirecmd

import (
	"fmt"

	bugout "github.com/bugout-dev/bugout-go/pkg"
	"github.com/spf13/cobra"
)

func CreateSpireCommand() *cobra.Command {
	spireCmd := &cobra.Command{
		Use:   "spire",
		Short: "Interact with Spire, the Bugout knowledge API",
		Long: `Bugout: The knowledge base for software teams

Spire is Bugout's knowledge API. You can use these commands to interact
with your personal and team knowledge bases from your command line.`,
	}

	pingCmd := CreatePingCommand()

	spireCmd.AddCommand(pingCmd)

	return spireCmd
}

func CreatePingCommand() *cobra.Command {
	pingCmd := &cobra.Command{
		Use:   "ping",
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
