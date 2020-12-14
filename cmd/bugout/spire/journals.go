package spirecmd

import (
	"encoding/json"

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

	journalsCmd.AddCommand(createCmd, deleteCmd, getCmd, listCmd, updateCmd)

	return journalsCmd
}

func CreateJournalsCreateCommand() *cobra.Command {
	var token, name string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Bugout journal",
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
	cmd.MarkFlagRequired("token")
	cmd.MarkFlagRequired("name")

	return cmd
}

func CreateJournalsDeleteCommand() *cobra.Command {
	var token, journalID string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a Bugout journal",
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
	cmd.Flags().StringVarP(&journalID, "id", "i", "", "ID of journal to delete")
	cmd.MarkFlagRequired("token")
	cmd.MarkFlagRequired("id")

	return cmd
}

func CreateJournalsGetCommand() *cobra.Command {
	var token, journalID string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a Bugout journal",
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
	cmd.Flags().StringVarP(&journalID, "id", "i", "", "ID of journal to delete")
	cmd.MarkFlagRequired("token")
	cmd.MarkFlagRequired("id")

	return cmd
}

func CreateJournalsListCommand() *cobra.Command {
	var token string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all Bugout journal accessible by a given user",
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
	cmd.MarkFlagRequired("token")

	return cmd
}

func CreateJournalsUpdateCommand() *cobra.Command {
	var token, journalID, name string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a Bugout journal",
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
	cmd.Flags().StringVarP(&journalID, "id", "i", "", "ID of journal to delete")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Name of journal to create")
	cmd.MarkFlagRequired("token")
	cmd.MarkFlagRequired("id")

	return cmd
}
