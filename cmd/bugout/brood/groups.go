package broodcmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bugout-dev/bugout-go/cmd/bugout/cmdutils"
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
	groupsRenameCmd := CreateGroupsRenameCommand()
	groupsAddUserCmd := CreateGroupsAddUserCommand()
	groupsRemoveUserCmd := CreateGroupsRemoveUserCommand()

	groupsCmd.AddCommand(groupsCreateCmd, groupsListCmd, groupsDeleteCmd, groupsRenameCmd, groupsAddUserCmd, groupsRemoveUserCmd)

	return groupsCmd
}

func CreateGroupsCreateCommand() *cobra.Command {
	var token, name string
	groupsCreateCmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a new bugout group",
		PreRunE: cmdutils.TokenArgPopulator,
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
	groupsCreateCmd.MarkFlagRequired("name")

	return groupsCreateCmd
}

func CreateGroupsListCommand() *cobra.Command {
	var token string
	groupsListCmd := &cobra.Command{
		Use:     "list",
		Short:   "List groups for a given user",
		PreRunE: cmdutils.TokenArgPopulator,
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

	return groupsListCmd
}

func CreateGroupsDeleteCommand() *cobra.Command {
	var token, groupID string
	groupsDeleteCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a bugout group (by ID)",
		PreRunE: cmdutils.TokenArgPopulator,
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
	groupsDeleteCmd.MarkFlagRequired("id")

	return groupsDeleteCmd
}

func CreateGroupsRenameCommand() *cobra.Command {
	var token, groupID, name string
	groupsRenameCmd := &cobra.Command{
		Use:     "rename",
		Short:   "Rename a bugout group",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			groups, groupsErr := client.Brood.RenameGroup(token, groupID, name)
			if groupsErr != nil {
				return groupsErr
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(groups)
			return encodeErr
		},
	}

	groupsRenameCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	groupsRenameCmd.Flags().StringVarP(&groupID, "id", "i", "", "ID of group to delete")
	groupsRenameCmd.Flags().StringVarP(&name, "name", "n", "", "Name of group to create")
	groupsRenameCmd.MarkFlagRequired("id")
	groupsRenameCmd.MarkFlagRequired("name")

	return groupsRenameCmd
}

func CreateGroupsAddUserCommand() *cobra.Command {
	var token, groupID, username, role string

	roles := []string{"owner", "member"}
	rolesStr := strings.Join(roles, ",")
	rolesHelp := fmt.Sprintf("User's role in group. Choices: %s", rolesStr)

	groupsAddUserCmd := &cobra.Command{
		Use:   "add-user",
		Short: "Add a Bugout user to a group",
		Args: func(cmd *cobra.Command, args []string) error {
			roleMatch := false
			for i := range roles {
				if role == roles[i] {
					roleMatch = true
					break
				}
			}

			if !roleMatch {
				return fmt.Errorf("Invalid role: %s. Choices: %s", role, rolesStr)
			}

			return nil
		},
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			membership, membershipErr := client.Brood.AddUserToGroup(token, groupID, username, role)
			if membershipErr != nil {
				return membershipErr
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(membership)
			return encodeErr
		},
	}

	groupsAddUserCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	groupsAddUserCmd.Flags().StringVarP(&groupID, "id", "i", "", "ID of group to add user to")
	groupsAddUserCmd.Flags().StringVarP(&username, "username", "u", "", "Bugout username of user to add to group")
	groupsAddUserCmd.Flags().StringVarP(&role, "role", "r", "", rolesHelp)
	groupsAddUserCmd.MarkFlagRequired("id")
	groupsAddUserCmd.MarkFlagRequired("username")
	groupsAddUserCmd.MarkFlagRequired("role")

	return groupsAddUserCmd
}

func CreateGroupsRemoveUserCommand() *cobra.Command {
	var token, groupID, username string

	groupsRemoveUserCmd := &cobra.Command{
		Use:     "remove-user",
		Short:   "Remove a Bugout user from a group",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			membership, membershipErr := client.Brood.RemoveUserFromGroup(token, groupID, username)
			if membershipErr != nil {
				return membershipErr
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(membership)
			return encodeErr
		},
	}

	groupsRemoveUserCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	groupsRemoveUserCmd.Flags().StringVarP(&groupID, "id", "i", "", "ID of group to add user to")
	groupsRemoveUserCmd.Flags().StringVarP(&username, "username", "u", "", "Bugout username of user to add to group")
	groupsRemoveUserCmd.MarkFlagRequired("id")
	groupsRemoveUserCmd.MarkFlagRequired("username")

	return groupsRemoveUserCmd
}
