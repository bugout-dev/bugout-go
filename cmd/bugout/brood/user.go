package broodcmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	bugout "github.com/bugout-dev/bugout-go/pkg"
)

func CreateUserCommand() *cobra.Command {
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "Bugout user operations",
	}

	userCreateCmd := CreateUserCreateCommand()
	userLoginCmd := CreateUserLoginCommand()

	userCmd.AddCommand(userCreateCmd, userLoginCmd)

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

func CreateUserLoginCommand() *cobra.Command {
	var username, password, tokenType, note string
	userLoginCmd := &cobra.Command{
		Use:   "login",
		Short: "Generate an access token for the given Bugout user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			token, tokenErr := client.Brood.GenerateToken(username, password)
			if tokenErr != nil {
				return tokenErr
			}

			if tokenType != "" || note != "" {
				_, annotationErr := client.Brood.AnnotateToken(token, tokenType, note)
				if annotationErr != nil {
					return annotationErr
				}
			}

			fmt.Println(token)

			return nil
		},
	}

	userLoginCmd.Flags().StringVarP(&username, "username", "u", "", "Desired username")
	userLoginCmd.Flags().StringVarP(&password, "password", "p", "", "Password for user")
	userLoginCmd.Flags().StringVarP(&tokenType, "type", "t", "", "Token type")
	userLoginCmd.Flags().StringVar(&note, "note", "Created using bugout CLI", "Note about the token")
	userLoginCmd.MarkFlagRequired("username")
	userLoginCmd.MarkFlagRequired("password")

	return userLoginCmd
}