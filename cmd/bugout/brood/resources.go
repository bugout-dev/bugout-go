package broodcmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bugout-dev/bugout-go/cmd/bugout/cmdutils"
	bugout "github.com/bugout-dev/bugout-go/pkg"
)

func GenerateResourcesCommand() *cobra.Command {
	resourcesCmd := &cobra.Command{
		Use:   "resources",
		Short: "Bugout resources operations",
	}

	resourcesCreateCmd := GenerateResourceCreateCommand()
	resourcesDeleteCmd := GenerateResourceDeleteCommand()
	resourcesGetCmd := GenerateResourcesGetCommand()

	resourcesCmd.AddCommand(resourcesCreateCmd, resourcesDeleteCmd, resourcesGetCmd)

	return resourcesCmd
}

func GenerateResourceCreateCommand() *cobra.Command {
	var token, applicationId string
	resourceCreateCmd := &cobra.Command{
		Use:     "create [resource]",
		Short:   "Create new resource",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("One resource argument as json string must be specified")
			}
			var resourceRaw interface{}
			err := json.Unmarshal([]byte(args[0]), &resourceRaw)
			if err != nil {
				return err
			}

			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			resource, err := client.Brood.CreateResource(token, applicationId, resourceRaw)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(&resource)
			return encodeErr
		},
	}

	resourceCreateCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	resourceCreateCmd.Flags().StringVarP(&applicationId, "application_id", "a", "", "Application ID resource belongs to")

	return resourceCreateCmd
}

func GenerateResourceDeleteCommand() *cobra.Command {
	var token, resourceId string
	resourceDeleteCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete resource",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			resource, err := client.Brood.DeleteResource(token, resourceId)
			if err != nil {
				return nil
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(&resource)
			return encodeErr
		},
	}

	resourceDeleteCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	resourceDeleteCmd.Flags().StringVarP(&resourceId, "resource_id", "r", "", "Resource ID")

	return resourceDeleteCmd
}

func GenerateResourcesGetCommand() *cobra.Command {
	var token, applicationId string
	var queryParams map[string]string
	resourcesGetCmd := &cobra.Command{
		Use:     "get",
		Short:   "Get resources of application",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			resources, err := client.Brood.GetResources(token, applicationId, queryParams)
			if err != nil {
				return nil
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(&resources)
			return encodeErr
		},
	}

	resourcesGetCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	resourcesGetCmd.Flags().StringVarP(&applicationId, "application_id", "a", "", "Application ID resource belongs to")
	resourcesGetCmd.Flags().StringToStringVarP(&queryParams, "params", "p", nil, "Optional query parameters to filter resources")

	return resourcesGetCmd
}
