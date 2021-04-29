package spirecmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bugout-dev/bugout-go/cmd/bugout/cmdutils"
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
	deleteCmd := CreateEntriesDeleteCommand()
	getCmd := CreateEntriesGetCommand()
	listCmd := CreateEntriesListCommand()
	searchCmd := CreateEntriesSearchCommand()
	tagCmd := CreateEntriesTagCommand()
	untagCmd := CreateEntriesUntagCommand()
	updateCmd := CreateEntriesUpdateCommand()
	cmd.AddCommand(createCmd, deleteCmd, getCmd, listCmd, searchCmd, tagCmd, untagCmd, updateCmd)

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
		PreRunE: cmdutils.CompositePopulator(cmdutils.TokenArgPopulator, cmdutils.JournalIDArgPopulator),
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
	cmd.MarkFlagRequired("title")
	cmd.MarkFlagFilename("file")

	return cmd
}

func CreateEntriesDeleteCommand() *cobra.Command {
	var token, journalID, entryID string
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete an entry from a Bugout journal",
		PreRunE: cmdutils.CompositePopulator(cmdutils.TokenArgPopulator, cmdutils.JournalIDArgPopulator),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			entry, err := client.Spire.DeleteEntry(token, journalID, entryID)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(entry)
			return encodeErr
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	cmd.Flags().StringVarP(&journalID, "journal", "j", "", "ID of journal")
	cmd.Flags().StringVarP(&entryID, "id", "i", "", "ID of entry")
	cmd.MarkFlagRequired("id")

	return cmd
}

func CreateEntriesGetCommand() *cobra.Command {
	var token, journalID, entryID string
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get an entry from a Bugout journal",
		PreRunE: cmdutils.CompositePopulator(cmdutils.TokenArgPopulator, cmdutils.JournalIDArgPopulator),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			entry, err := client.Spire.GetEntry(token, journalID, entryID)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(entry)
			return encodeErr
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	cmd.Flags().StringVarP(&journalID, "journal", "j", "", "ID of journal")
	cmd.Flags().StringVarP(&entryID, "id", "i", "", "ID of entry")
	cmd.MarkFlagRequired("id")

	return cmd
}

func CreateEntriesListCommand() *cobra.Command {
	var token, journalID string
	var limit, offset int
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all entries in a Bugout journal",
		PreRunE: cmdutils.CompositePopulator(cmdutils.TokenArgPopulator, cmdutils.JournalIDArgPopulator),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			entries, err := client.Spire.ListEntries(token, journalID, limit, offset)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(entries)
			return encodeErr
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	cmd.Flags().StringVarP(&journalID, "journal", "j", "", "ID of journal")
	cmd.Flags().IntVarP(&limit, "limit", "N", 10, "Number of entries per page")
	cmd.Flags().IntVarP(&offset, "offset", "n", 0, "Index of starting entry on current page")

	return cmd
}

func CreateEntriesSearchCommand() *cobra.Command {
	var token, journalID string
	var limit, offset int
	var queryParams map[string]string
	cmd := &cobra.Command{
		Use:     "search [query]",
		Short:   "Search across the entries in a Bugout journal",
		PreRunE: cmdutils.CompositePopulator(cmdutils.TokenArgPopulator, cmdutils.JournalIDArgPopulator),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			searchQuery := strings.Join(args, " ")

			entries, err := client.Spire.SearchEntries(token, journalID, searchQuery, limit, offset, queryParams)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(entries)
			return encodeErr
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	cmd.Flags().StringVarP(&journalID, "journal", "j", "", "ID of journal")
	cmd.Flags().IntVarP(&limit, "limit", "N", 10, "Number of entries per page")
	cmd.Flags().IntVarP(&offset, "offset", "n", 0, "Index of starting entry on current page")
	cmd.Flags().StringToStringVarP(&queryParams, "params", "p", nil, "Optional query parameters to add to the query")

	return cmd
}

func CreateEntriesTagCommand() *cobra.Command {
	var token, journalID, entryID string
	cmd := &cobra.Command{
		Use:     "tag [tags...]",
		Short:   "Add tags to an entry",
		PreRunE: cmdutils.CompositePopulator(cmdutils.TokenArgPopulator, cmdutils.JournalIDArgPopulator),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			entry, err := client.Spire.TagEntry(token, journalID, entryID, args)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(entry)
			return encodeErr
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	cmd.Flags().StringVarP(&journalID, "journal", "j", "", "ID of journal")
	cmd.Flags().StringVarP(&entryID, "id", "i", "", "ID of entry")
	cmd.MarkFlagRequired("id")

	return cmd
}

func CreateEntriesUntagCommand() *cobra.Command {
	var token, journalID, entryID string
	cmd := &cobra.Command{
		Use:     "untag [tags...]",
		Short:   "Remove tags from an entry",
		PreRunE: cmdutils.CompositePopulator(cmdutils.TokenArgPopulator, cmdutils.JournalIDArgPopulator),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
			}

			entry, err := client.Spire.UntagEntry(token, journalID, entryID, args)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(entry)
			return encodeErr
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	cmd.Flags().StringVarP(&journalID, "journal", "j", "", "ID of journal")
	cmd.Flags().StringVarP(&entryID, "id", "i", "", "ID of entry")
	cmd.MarkFlagRequired("id")

	return cmd
}

func CreateEntriesUpdateCommand() *cobra.Command {
	var token, journalID, entryID, title, content, contentFile string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an entry in a Bugout journal",
		Args: func(cmd *cobra.Command, args []string) error {
			if (content == "" && contentFile == "") || (content != "" && contentFile != "") {
				return errors.New("Exactly one of --content or --file must be specified")
			}

			return nil
		},
		PreRunE: cmdutils.CompositePopulator(cmdutils.TokenArgPopulator, cmdutils.JournalIDArgPopulator),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := bugout.ClientFromEnv()
			if clientErr != nil {
				return clientErr
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

			entry, err := client.Spire.UpdateEntry(token, journalID, entryID, title, entryContent)
			if err != nil {
				return err
			}

			encodeErr := json.NewEncoder(cmd.OutOrStdout()).Encode(entry)
			return encodeErr
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "Bugout access token to use for the request")
	cmd.Flags().StringVarP(&journalID, "journal", "j", "", "ID of journal")
	cmd.Flags().StringVarP(&entryID, "id", "i", "", "ID of entry")
	cmd.Flags().StringVar(&title, "title", "", "Title of new entry")
	cmd.Flags().StringVarP(&content, "content", "c", "", "Content of entry")
	cmd.Flags().StringVarP(&contentFile, "file", "f", "", "File containing contents of entry")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("title")
	cmd.MarkFlagFilename("file")

	return cmd
}
