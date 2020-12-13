package broodcmd

import (
	"fmt"

	"github.com/spf13/cobra"

	bugout "github.com/bugout-dev/bugout-go/pkg"
)

func CreateBroodCommand() *cobra.Command {
	broodCmd := &cobra.Command{
		Use:   "brood",
		Short: "Interact with Brood, the Bugout authentication API, from your command line",
		Long: `Bugout: The knowledge base for software teams

Brood is Bugout's authentication API. You can use these commands to interact with Bugout users and
groups from your command line`,
	}

	pingCmd := CreatePingCommand()
	userCmd := CreateUserCommand()
	versionCmd := CreateVersionCommand()

	broodCmd.AddCommand(pingCmd, userCmd, versionCmd)

	return broodCmd
}

func CreatePingCommand() *cobra.Command {
	pingCmd := &cobra.Command{
		Use:   "ping",
		Short: "Ping Brood to see if it is active",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			result, err := client.Brood.Ping()
			if err != nil {
				return err
			}

			fmt.Println(result)

			return nil
		},
	}

	return pingCmd
}

func CreateVersionCommand() *cobra.Command {
	pingCmd := &cobra.Command{
		Use:   "version",
		Short: "Check Brood version",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			result, err := client.Brood.Version()
			if err != nil {
				return err
			}

			fmt.Println(result)

			return nil
		},
	}

	return pingCmd
}
