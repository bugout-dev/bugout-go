package broodcmd

import (
	"encoding/json"

	"github.com/spf13/cobra"

	bugout "github.com/bugout-dev/bugout-go/pkg"
)

func CreateGroupsCommand() *cobra.Command {
	groupsCmd := &cobra.Command{
		Use:   "groups",
		Short: "Bugout group operations",
	}

	groupsCreateCmd := CreateGroupsCreateCommand()
	groupsListCmd := CreateGroupsListCommand()
	groupsDeleteCmd := CreateGroupsDeleteCommand()

	groupsCmd.AddCommand(groupsCreateCmd, groupsListCmd, groupsDeleteCmd)

	return groupsCmd
}

func CreateGroupsCreateCommand() *cobra.Command {
	var token, name string
	groupsCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new bugout groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			groups, groupsErr := client.Brood.CreateGroup(token, name)
			if groupsErr != nil {
				return groupsErr
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(groups)
			return encodeErr
		},
	}

	groupsCreateCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	groupsCreateCmd.Flags().StringVarP(&name, "name", "n", "", "Name of group to create")
	groupsCreateCmd.MarkFlagRequired("token")
	groupsCreateCmd.MarkFlagRequired("name")

	return groupsCreateCmd
}

func CreateGroupsListCommand() *cobra.Command {
	var token string
	groupsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List groups for a given user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			groups, groupsErr := client.Brood.GetUserGroups(token)
			if groupsErr != nil {
				return groupsErr
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(groups)
			return encodeErr
		},
	}

	groupsListCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	groupsListCmd.MarkFlagRequired("token")

	return groupsListCmd
}

func CreateGroupsDeleteCommand() *cobra.Command {
	var token, groupID string
	groupsDeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a bugout group (by ID)",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			group, groupErr := client.Brood.DeleteGroup(token, groupID)
			if groupErr != nil {
				return groupErr
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(group)
			return encodeErr
		},
	}

	groupsDeleteCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	groupsDeleteCmd.Flags().StringVarP(&groupID, "id", "i", "", "ID of group to delete")
	groupsDeleteCmd.MarkFlagRequired("token")
	groupsDeleteCmd.MarkFlagRequired("id")

	return groupsDeleteCmd
}
