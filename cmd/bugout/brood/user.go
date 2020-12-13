package broodcmd

import (
	"encoding/json"

	"github.com/spf13/cobra"

	bugout "github.com/bugout-dev/bugout-go/pkg"
)

func CreateUserCommand() *cobra.Command {
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "Bugout user operations",
	}

	userCreateCmd := CreateUserCreateCommand()

	userCmd.AddCommand(userCreateCmd)

	return userCmd
}

func CreateUserCreateCommand() *cobra.Command {
	var username, email, password string
	userCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new bugout user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			user, userErr := client.Brood.CreateUser(username, email, password)
			if userErr != nil {
				return userErr
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(user)
			return encodeErr
		},
	}

	userCreateCmd.Flags().StringVarP(&username, "username", "u", "", "Desired username")
	userCreateCmd.Flags().StringVarP(&email, "email", "e", "", "Email address for user")
	userCreateCmd.Flags().StringVarP(&password, "password", "p", "", "Password for user")
	userCreateCmd.MarkFlagRequired("username")
	userCreateCmd.MarkFlagRequired("email")
	userCreateCmd.MarkFlagRequired("password")

	return userCreateCmd
}
