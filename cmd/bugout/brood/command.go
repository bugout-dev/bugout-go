package broodcmd

import (
	"fmt"

	"github.com/spf13/cobra"

	bugout "github.com/bugout-dev/bugout-go/pkg"
)

func PopulateBroodCommands(cmd *cobra.Command) {
	groupsCmd := CreateGroupsCommand()
	resourcesCmd := GenerateResourcesCommand()
	pingCmd := CreatePingCommand()
	userCmd := CreateUserCommand()
	versionCmd := CreateVersionCommand()

	cmd.AddCommand(groupsCmd, resourcesCmd, pingCmd, userCmd, versionCmd)
}

func CreatePingCommand() *cobra.Command {
	pingCmd := &cobra.Command{
		Use:   "ping-brood",
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
		Use:   "version-brood",
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
