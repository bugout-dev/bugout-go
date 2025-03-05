package broodcmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bugout-dev/bugout-go/cmd/bugout/cmdutils"
	bugout "github.com/bugout-dev/bugout-go/pkg"
)

func CreateUserCommand() *cobra.Command {
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "Bugout user operations",
	}

	userAuthCmd := CreateUserAuthCommand()
	userCreateCmd := CreateUserCreateCommand()
	userLoginCmd := CreateUserLoginCommand()
	userTokensCmd := CreateUserTokensCommand()
	userGetCmd := CreateUserGetCommand()
	userFindCmd := CreateUserFindCommand()
	userVerifyCmd := CreateUserVerifyCommand()
	userChangePasswordCmd := CreateUserChangePasswordCommand()

	userCmd.AddCommand(userAuthCmd, userCreateCmd, userLoginCmd, userTokensCmd, userGetCmd, userFindCmd, userVerifyCmd, userChangePasswordCmd)

	return userCmd
}

func CreateUserAuthCommand() *cobra.Command {
	var token string
	userGetCmd := &cobra.Command{
		Use:     "auth",
		Short:   "Get the user with groups represented by a token",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			userAuth, err := client.Brood.Auth(token)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(&userAuth)
			return encodeErr
		},
	}

	userGetCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")

	return userGetCmd
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

func CreateUserTokensCommand() *cobra.Command {
	var token string
	userTokensCmd := &cobra.Command{
		Use:   "tokens",
		Short: "List all access tokens for a given user",
		Long: `List all access tokens for a given user.

The user is identified by a Bugout access token.`,
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			tokens, err := client.Brood.ListTokens(token)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(&tokens)
			return encodeErr
		},
	}

	userTokensCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")

	return userTokensCmd
}

func CreateUserGetCommand() *cobra.Command {
	var token string
	userGetCmd := &cobra.Command{
		Use:     "get",
		Short:   "Get the user represented by a token",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			user, err := client.Brood.GetUser(token)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(&user)
			return encodeErr
		},
	}

	userGetCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")

	return userGetCmd
}

func CreateUserFindCommand() *cobra.Command {
	var token, user_id, username, email, application_id string
	userFindCmd := &cobra.Command{
		Use:   "find",
		Short: "Find user if exists",
		Args: func(cmd *cobra.Command, args []string) error {
			if user_id == "" && username == "" && email == "" && application_id == "" {
				return errors.New("Exactly one of --user_id or --username or --email or application_id must be specified")
			}

			return nil
		},
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			queryParams := make(map[string]string)
			if user_id != "" {
				queryParams["user_id"] = user_id
			}
			if username != "" {
				queryParams["username"] = username
			}
			if email != "" {
				queryParams["email"] = email
			}
			if application_id != "" {
				queryParams["application_id"] = application_id
			}

			user, err := client.Brood.FindUser(token, queryParams)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(&user)
			return encodeErr
		},
	}

	userFindCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	userFindCmd.Flags().StringVarP(&user_id, "user_id", "i", "", "Bugout user ID")
	userFindCmd.Flags().StringVarP(&username, "username", "u", "", "User name")
	userFindCmd.Flags().StringVarP(&email, "email", "e", "", "User email address")
	userFindCmd.Flags().StringVarP(&application_id, "application_id", "a", "", "Application user belongs to")

	return userFindCmd
}

func CreateUserVerifyCommand() *cobra.Command {
	var token, code string
	userVerifyCmd := &cobra.Command{
		Use:     "verify",
		Short:   "Verify the user represented by a token",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			user, err := client.Brood.VerifyUser(token, code)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(&user)
			return encodeErr
		},
	}

	userVerifyCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	userVerifyCmd.Flags().StringVarP(&code, "code", "c", "", "Verification code that was sent to user's email address")
	userVerifyCmd.MarkFlagRequired("code")

	return userVerifyCmd
}

func CreateUserChangePasswordCommand() *cobra.Command {
	var token, currentPassword, newPassword string
	userChangePasswordCmd := &cobra.Command{
		Use:     "change-password",
		Short:   "ChangePassword the user represented by a token",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			user, err := client.Brood.ChangePassword(token, currentPassword, newPassword)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(&user)
			return encodeErr
		},
	}

	userChangePasswordCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	userChangePasswordCmd.Flags().StringVar(&currentPassword, "current", "", "Current password for the given user")
	userChangePasswordCmd.Flags().StringVar(&newPassword, "new", "", "New password for the given user")
	userChangePasswordCmd.MarkFlagRequired("current")
	userChangePasswordCmd.MarkFlagRequired("new")

	return userChangePasswordCmd
}
