package spirecmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	bugout "github.com/bugout-dev/bugout-go/pkg"
	"github.com/bugout-dev/bugout-go/pkg/spire"
	"github.com/spf13/cobra"
)

func CreateEntriesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entries",
		Short: "Interact with entries in Bugout journals",
	}

	createCmd := CreateEntriesCreateCommand()
	cmd.AddCommand(createCmd)

	return cmd
}

func CreateEntriesCreateCommand() *cobra.Command {
	var token, journalID, title, content, contentFile string
	var tags []string
	var contextMap map[string]string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new entry in a Bugout journal",
		Args: func(cmd *cobra.Command, args []string) error {
			validContextKeys := map[string]bool{
				"type": true,
				"id":   true,
				"url":  true,
			}
			for k := range contextMap {
				if _, exists := validContextKeys[k]; !exists {
					return fmt.Errorf("Invalid key (%s) in context. Valid choices are: type,id,url", k)
				}
			}

			if (content == "" && contentFile == "") || (content != "" && contentFile != "") {
				return errors.New("Exactly one of --content or --file must be specified")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			context := spire.EntryContext{}
			if value, exists := contextMap["type"]; exists {
				context.ContextType = value
			}
			if value, exists := contextMap["id"]; exists {
				context.ContextID = value
			}
			if value, exists := contextMap["url"]; exists {
				context.ContextURL = value
			}

			var entryContent string
			if content != "" {
				entryContent = content
			} else {
				contentFileBytes, readErr := ioutil.ReadFile(contentFile)
				if readErr != nil {
					return readErr
				}
				entryContent = string(contentFileBytes)
			}

			entry, err := client.Spire.CreateEntry(token, journalID, title, entryContent, tags, context)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(entry)
			return encodeErr
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	cmd.Flags().StringVarP(&journalID, "journal", "j", "", "ID of journal")
	cmd.Flags().StringVar(&title, "title", "", "Title of new entry")
	cmd.Flags().StringVarP(&content, "content", "c", "", "Content of entry")
	cmd.Flags().StringVarP(&contentFile, "file", "f", "", "File containing contents of entry")
	cmd.Flags().StringSliceVar(&tags, "tags", []string{}, "Tags to apply to the new entry (as a comma-separated list of strings)")
	cmd.Flags().StringToStringVar(&contextMap, "context", map[string]string{}, "Context for the new entry (in the format type=<context type>,id=<context id>,url=<context url>)")
	cmd.MarkFlagRequired("token")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("title")
	cmd.MarkFlagFilename("file")

	return cmd
}
