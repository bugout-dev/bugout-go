package broodcmd

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/bugout-dev/bugout-go/cmd/bugout/cmdutils"
	bugout "github.com/bugout-dev/bugout-go/pkg"
)

func CreateApplicationsCommand() *cobra.Command {
	applicationsCmd := &cobra.Command{
		Use:   "applications",
		Short: "Bugout application operations",
	}

	applicationsCreateCmd := CreateApplicationsCreateCommand()
	applicationsGetCmd := CreateApplicationsGetCommand()
	applicationsListCmd := CreateApplicationsListCommand()
	applicationsDeleteCmd := CreateApplicationsDeleteCommand()

	applicationsCmd.AddCommand(applicationsCreateCmd, applicationsGetCmd, applicationsListCmd, applicationsDeleteCmd)

	return applicationsCmd
}

func CreateApplicationsCreateCommand() *cobra.Command {
	var token, groupId, name, description string
	applicationsCreateCmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a new bugout application",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			applications, applicationsErr := client.Brood.CreateApplication(token, groupId, name, description)
			if applicationsErr != nil {
				return applicationsErr
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(applications)
			return encodeErr
		},
	}

	applicationsCreateCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	applicationsCreateCmd.Flags().StringVarP(&groupId, "group", "g", "", "ID of Brood group to create the application under")
	applicationsCreateCmd.Flags().StringVarP(&name, "name", "n", "", "Name of application to create")
	applicationsCreateCmd.Flags().StringVarP(&description, "description", "d", "", "Description of application to create")
	applicationsCreateCmd.MarkFlagRequired("group")
	applicationsCreateCmd.MarkFlagRequired("name")

	return applicationsCreateCmd
}

func CreateApplicationsGetCommand() *cobra.Command {
	var token, applicationId string
	applicationsGetCmd := &cobra.Command{
		Use:     "get",
		Short:   "Get an application by ID",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			applications, applicationsErr := client.Brood.GetApplication(token, applicationId)
			if applicationsErr != nil {
				return applicationsErr
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(applications)
			return encodeErr
		},
	}

	applicationsGetCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	applicationsGetCmd.Flags().StringVarP(&applicationId, "application", "a", "", "ID of application to get")
	applicationsGetCmd.MarkFlagRequired("application")

	return applicationsGetCmd
}

func CreateApplicationsListCommand() *cobra.Command {
	var token, groupId string
	applicationsListCmd := &cobra.Command{
		Use:     "list",
		Short:   "List applications for a given user",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			applications, applicationsErr := client.Brood.ListApplications(token, groupId)
			if applicationsErr != nil {
				return applicationsErr
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(applications)
			return encodeErr
		},
	}

	applicationsListCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	applicationsListCmd.Flags().StringVarP(&groupId, "group", "g", "", "Only return applications owned by this group (optional)")

	return applicationsListCmd
}

func CreateApplicationsDeleteCommand() *cobra.Command {
	var token, applicationID string
	applicationsDeleteCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a bugout application (by ID)",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bugout.ClientFromEnv()
			if err != nil {
				return err
			}

			application, applicationErr := client.Brood.DeleteApplication(token, applicationID)
			if applicationErr != nil {
				return applicationErr
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(application)
			return encodeErr
		},
	}

	applicationsDeleteCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	applicationsDeleteCmd.Flags().StringVarP(&applicationID, "id", "i", "", "ID of application to delete")
	applicationsDeleteCmd.MarkFlagRequired("id")

	return applicationsDeleteCmd
}
