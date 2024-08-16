package broodcmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bugout-dev/bugout-go/cmd/bugout/cmdutils"
	bugout "github.com/bugout-dev/bugout-go/pkg"
	"github.com/bugout-dev/bugout-go/pkg/brood"
)

func GenerateResourcesCommand() *cobra.Command {
	resourcesCmd := &cobra.Command{
		Use:   "resources",
		Short: "Bugout resources operations",
	}

	resourcesCreateCmd := GenerateResourceCreateCommand()
	resourcesUpdateCmd := GenerateResourceUpdateCommand()
	resourcesDeleteCmd := GenerateResourceDeleteCommand()
	resourceGetCmd := GenerateResourceGetCommand()
	resourcesListCmd := GenerateResourcesListCommand()

	resourceHoldersCmd := GenerateResourceHoldersCommand()

	resourcesCmd.AddCommand(resourcesCreateCmd, resourcesUpdateCmd, resourcesDeleteCmd, resourceGetCmd, resourcesListCmd, resourceHoldersCmd)

	return resourcesCmd
}

func GenerateResourceHoldersCommand() *cobra.Command {
	holdersCmd := &cobra.Command{
		Use:   "holders",
		Short: "Bugout resource holders operations",
	}

	resourceHoldersGetCmd := GenerateResourceHoldersGetCommand()
	resourceHoldersAddCmd := GenerateResourceHoldersAddCommand()
	resourceHoldersDeleteCmd := GenerateResourceHoldersDeleteCommand()

	holdersCmd.AddCommand(resourceHoldersGetCmd, resourceHoldersAddCmd, resourceHoldersDeleteCmd)

	return holdersCmd
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

type resourceUpdateArg struct {
	Update   interface{} `json:"update"`
	DropKeys []string    `json:"drop_keys"`
}

func GenerateResourceUpdateCommand() *cobra.Command {
	var token, resourceId string
	resourceUpdateCmd := &cobra.Command{
		Use:     "update [resource]",
		Short:   "Update resource",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf(`Resource argument as json string in format {"update": {"key": "value"}, "drop_keys": ["key"]} must be specified`)
			}
			var resourceUpdateRaw resourceUpdateArg
			err := json.Unmarshal([]byte(args[0]), &resourceUpdateRaw)
			if err != nil {
				return err
			}

			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			resource, err := client.Brood.UpdateResource(token, resourceId, resourceUpdateRaw.Update, resourceUpdateRaw.DropKeys)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(&resource)
			return encodeErr
		},
	}

	resourceUpdateCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	resourceUpdateCmd.Flags().StringVarP(&resourceId, "resource_id", "r", "", "Resource ID")

	return resourceUpdateCmd
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

func GenerateResourceGetCommand() *cobra.Command {
	var token, resourceId string
	resourcesGetCmd := &cobra.Command{
		Use:     "get",
		Short:   "Get resource of application",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			resources, err := client.Brood.GetResource(token, resourceId)
			if err != nil {
				return nil
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(&resources)
			return encodeErr
		},
	}

	resourcesGetCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	resourcesGetCmd.Flags().StringVarP(&resourceId, "resource_id", "r", "", "Resource ID")

	return resourcesGetCmd
}

func GenerateResourcesListCommand() *cobra.Command {
	var token, applicationId string
	var queryParams map[string]string
	resourcesGetCmd := &cobra.Command{
		Use:     "list",
		Short:   "List resources of application",
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

func GenerateResourceHoldersGetCommand() *cobra.Command {
	var token, resourceId string
	resourceHoldersGetCmd := &cobra.Command{
		Use:     "get",
		Short:   "Get resource holders",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			resourceHolders, err := client.Brood.GetResourceHolders(token, resourceId)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(&resourceHolders)
			return encodeErr
		},
	}

	resourceHoldersGetCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	resourceHoldersGetCmd.Flags().StringVarP(&resourceId, "resource_id", "r", "", "Resource ID")

	return resourceHoldersGetCmd
}

func GenerateResourceHoldersAddCommand() *cobra.Command {
	var token, resourceId string
	resourceHoldersAddCmd := &cobra.Command{
		Use:     "add [holder]",
		Short:   "Add resource holders",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf(`Resource holder argument as json string in format {"holder_id": "user_or_group_uuid", "holder_type": "user_or_group", "permissions": ["admin", "create", "read", "update", "delete"]} must be specified`)
			}
			var resourceHolderPermissionsRaw brood.ResourceHolder
			err := json.Unmarshal([]byte(args[0]), &resourceHolderPermissionsRaw)
			if err != nil {
				return err
			}

			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			resource, err := client.Brood.AddResourceHolderPermissions(token, resourceId, resourceHolderPermissionsRaw)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(&resource)
			return encodeErr
		},
	}

	resourceHoldersAddCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	resourceHoldersAddCmd.Flags().StringVarP(&resourceId, "resource_id", "r", "", "Resource ID")

	return resourceHoldersAddCmd
}

func GenerateResourceHoldersDeleteCommand() *cobra.Command {
	var token, resourceId string
	resourceHoldersDeleteCmd := &cobra.Command{
		Use:     "delete [holder]",
		Short:   "Delete resource holders",
		PreRunE: cmdutils.TokenArgPopulator,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf(`Resource holder argument as json string in format {"holder_id": "user_or_group_uuid", "holder_type": "user_or_group", "permissions": ["admin", "create", "read", "update", "delete"]} must be specified`)
			}
			var resourceHolderPermissionsRaw brood.ResourceHolder
			err := json.Unmarshal([]byte(args[0]), &resourceHolderPermissionsRaw)
			if err != nil {
				return err
			}

			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			resource, err := client.Brood.DeleteResourceHolderPermissions(token, resourceId, resourceHolderPermissionsRaw)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(&resource)
			return encodeErr
		},
	}

	resourceHoldersDeleteCmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	resourceHoldersDeleteCmd.Flags().StringVarP(&resourceId, "resource_id", "r", "", "Resource ID")

	return resourceHoldersDeleteCmd
}
