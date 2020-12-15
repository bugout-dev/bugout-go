package spirecmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bugout-dev/bugout-go/cmd/bugout/cmdutils"
	bugout "github.com/bugout-dev/bugout-go/pkg"
	"github.com/bugout-dev/bugout-go/pkg/spire"
	"github.com/spf13/cobra"
)

func CreateJournalsCommand() *cobra.Command {
	journalsCmd := &cobra.Command{
		Use:   "journals",
		Short: "Interact with Bugout journals",
	}

	createCmd := CreateJournalsCreateCommand()
	deleteCmd := CreateJournalsDeleteCommand()
	getCmd := CreateJournalsGetCommand()
	listCmd := CreateJournalsListCommand()
	updateCmd := CreateJournalsUpdateCommand()
	addMemberCmd := CreateJournalsAddMemberCommand()
	removeMemberCmd := CreateJournalsRemoveMemberCommand()

	journalsCmd.AddCommand(createCmd, deleteCmd, getCmd, listCmd, updateCmd, addMemberCmd, removeMemberCmd)

	return journalsCmd
}

func CreateJournalsCreateCommand() *cobra.Command {
	var token, name string
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a new Bugout journal",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			journal, err := client.Spire.CreateJournal(token, name)
			if err != nil {
				return err
			}
			json.NewEncoder(cmd.OutOrStdout()).Encode(journal)
			return nil
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Name of journal to create")
	cmd.MarkFlagRequired("name")

	return cmd
}

func CreateJournalsDeleteCommand() *cobra.Command {
	var token, journalID string
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a Bugout journal",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			journal, err := client.Spire.DeleteJournal(token, journalID)
			if err != nil {
				return err
			}
			json.NewEncoder(cmd.OutOrStdout()).Encode(journal)
			return nil
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	cmd.Flags().StringVarP(&journalID, "journal", "j", "", "ID of journal to delete")
	cmd.MarkFlagRequired("journal")

	return cmd
}

func CreateJournalsGetCommand() *cobra.Command {
	var token, journalID string
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get a Bugout journal",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			journal, err := client.Spire.GetJournal(token, journalID)
			if err != nil {
				return err
			}
			json.NewEncoder(cmd.OutOrStdout()).Encode(journal)
			return nil
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	cmd.Flags().StringVarP(&journalID, "journal", "j", "", "ID of journal to get")
	cmd.MarkFlagRequired("journal")

	return cmd
}

func CreateJournalsListCommand() *cobra.Command {
	var token string
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all Bugout journal accessible by a given user",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			journals, err := client.Spire.ListJournals(token)
			if err != nil {
				return err
			}
			json.NewEncoder(cmd.OutOrStdout()).Encode(journals)
			return nil
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")

	return cmd
}

func CreateJournalsUpdateCommand() *cobra.Command {
	var token, journalID, name string
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update a Bugout journal",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			var journal spire.Journal
			var err error
			if name != "" {
				journal, err = client.Spire.UpdateJournal(token, journalID, name)
			} else {
				journal, err = client.Spire.GetJournal(token, journalID)
			}

			if err != nil {
				return err
			}

			json.NewEncoder(cmd.OutOrStdout()).Encode(journal)
			return nil
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	cmd.Flags().StringVarP(&journalID, "journal", "j", "", "ID of journal to update")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Updated name for journal")
	cmd.MarkFlagRequired("journal")

	return cmd
}

func CreateJournalsAddMemberCommand() *cobra.Command {
	var token, journalID, memberID, memberType string
	cmd := &cobra.Command{
		Use:     "add-member [permissions...]",
		Short:   "Add a member to a Bugout journal.",
		Long:    fmt.Sprintf("Add a member to a Bugout journal.\n\nValid permissions: %s", strings.Join(spire.ValidJournalPermissions(), ",")),
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			permissionsList, err := client.Spire.AddJournalMember(token, journalID, memberID, memberType, args)
			if err != nil {
				return err
			}

			json.NewEncoder(cmd.OutOrStdout()).Encode(permissionsList)
			return nil
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	cmd.Flags().StringVarP(&journalID, "journal", "j", "", "ID of journal to add a member to")
	cmd.Flags().StringVar(&memberID, "member", "", "ID for user or group to add as a member")
	cmd.Flags().StringVar(&memberType, "member-type", "user", fmt.Sprintf("Type of member (choices: %s)", strings.Join(spire.ValidMemberTypes(), ",")))
	cmd.MarkFlagRequired("journal")
	cmd.MarkFlagRequired("member")

	return cmd
}

func CreateJournalsRemoveMemberCommand() *cobra.Command {
	var token, journalID, memberID, memberType string
	cmd := &cobra.Command{
		Use:     "remove-member [permissions...]",
		Short:   "Remove a member from a Bugout journal.",
		Long:    fmt.Sprintf("Remove a member's permissions to a Bugout journal.\n\nValid permissions: %s", strings.Join(spire.ValidJournalPermissions(), ",")),
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			permissionsList, err := client.Spire.RemoveJournalMember(token, journalID, memberID, memberType, args)
			if err != nil {
				return err
			}

			json.NewEncoder(cmd.OutOrStdout()).Encode(permissionsList)
			return nil
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	cmd.Flags().StringVarP(&journalID, "journal", "j", "", "ID of journal to add a member to")
	cmd.Flags().StringVar(&memberID, "member", "", "ID for user or group to add as a member")
	cmd.Flags().StringVar(&memberType, "member-type", "user", fmt.Sprintf("Type of member (choices: %s)", strings.Join(spire.ValidMemberTypes(), ",")))
	cmd.MarkFlagRequired("journal")
	cmd.MarkFlagRequired("member")

	return cmd
}
