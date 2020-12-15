package spire

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bugout-dev/bugout-go/pkg/utils"
)

func (client SpireClient) CreateEntry(token, journalID, title, content string, tags []string, context EntryContext) (Entry, error) {
	entriesRoute := fmt.Sprintf("%s/%s/entries", client.Routes.Journals, journalID)
	requestBody := entryCreateRequest{
		Title:   title,
		Content: content,
		Tags:    tags,
	}

	if context.ContextType != "" {
		requestBody.ContextType = context.ContextType
	}
	if context.ContextID != "" {
		requestBody.ContextID = context.ContextID
	}
	if context.ContextURL != "" {
		requestBody.ContextURL = context.ContextURL
	}

	requestBuffer := new(bytes.Buffer)
	encodeErr := json.NewEncoder(requestBuffer).Encode(requestBody)
	if encodeErr != nil {
		return Entry{}, encodeErr
	}
	request, requestErr := http.NewRequest("POST", entriesRoute, requestBuffer)
	if requestErr != nil {
		return Entry{}, requestErr
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	response, responseErr := client.HTTPClient.Do(request)
	if responseErr != nil {
		return Entry{}, responseErr
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return Entry{}, statusErr
	}

	var entry Entry
	decodeErr := json.NewDecoder(response.Body).Decode(&entry)
	return entry, decodeErr
}

func (client SpireClient) DeleteEntry(token, journalID, entryID string) (Entry, error) {
	entryRoute := fmt.Sprintf("%s/%s/entries/%s", client.Routes.Journals, journalID, entryID)
	request, requestErr := http.NewRequest("DELETE", entryRoute, nil)
	if requestErr != nil {
		return Entry{}, requestErr
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	response, responseErr := client.HTTPClient.Do(request)
	if responseErr != nil {
		return Entry{}, responseErr
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return Entry{}, statusErr
	}

	var entry Entry
	decodeErr := json.NewDecoder(response.Body).Decode(&entry)
	return entry, decodeErr
}

func (client SpireClient) GetEntry(token, journalID, entryID string) (Entry, error) {
	entryRoute := fmt.Sprintf("%s/%s/entries/%s", client.Routes.Journals, journalID, entryID)
	request, requestErr := http.NewRequest("GET", entryRoute, nil)
	if requestErr != nil {
		return Entry{}, requestErr
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	response, responseErr := client.HTTPClient.Do(request)
	if responseErr != nil {
		return Entry{}, responseErr
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return Entry{}, statusErr
	}

	var entry Entry
	decodeErr := json.NewDecoder(response.Body).Decode(&entry)
	return entry, decodeErr
}

func (client SpireClient) ListEntries(token, journalID string, limit, offset int) (EntryResultsPage, error) {
	entriesRoute := fmt.Sprintf("%s/%s/search", client.Routes.Journals, journalID)
	request, requestErr := http.NewRequest("GET", entriesRoute, nil)
	if requestErr != nil {
		return EntryResultsPage{}, requestErr
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	// Pattern taken from Stack Overflow: https://stackoverflow.com/a/30657518/13659585
	query := request.URL.Query()
	query.Add("q", "")
	query.Add("limit", strconv.Itoa(limit))
	query.Add("offset", strconv.Itoa(offset))
	request.URL.RawQuery = query.Encode()

	response, responseErr := client.HTTPClient.Do(request)
	if responseErr != nil {
		return EntryResultsPage{}, responseErr
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return EntryResultsPage{}, statusErr
	}

	var entries EntryResultsPage
	decodeErr := json.NewDecoder(response.Body).Decode(&entries)
	return entries, decodeErr
}

func (client SpireClient) SearchEntries(token, journalID, searchQuery string, limit, offset int) (EntryResultsPage, error) {
	entriesRoute := fmt.Sprintf("%s/%s/search", client.Routes.Journals, journalID)
	request, requestErr := http.NewRequest("GET", entriesRoute, nil)
	if requestErr != nil {
		return EntryResultsPage{}, requestErr
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	// Pattern taken from Stack Overflow: https://stackoverflow.com/a/30657518/13659585
	query := request.URL.Query()
	query.Add("q", searchQuery)
	query.Add("limit", strconv.Itoa(limit))
	query.Add("offset", strconv.Itoa(offset))
	request.URL.RawQuery = query.Encode()

	response, responseErr := client.HTTPClient.Do(request)
	if responseErr != nil {
		return EntryResultsPage{}, responseErr
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return EntryResultsPage{}, statusErr
	}

	var entries EntryResultsPage
	decodeErr := json.NewDecoder(response.Body).Decode(&entries)
	return entries, decodeErr
}

func (client SpireClient) TagEntry(token, journalID, entryID string, tags []string) (Entry, error) {
	// The tagging endpoint returns an error if we try to add tags to an entry that it already has.
	// Therefore, we must first restrict the slice of new tags to only those tags which are not
	// already present on the entry.
	currentEntry, currentEntryErr := client.GetEntry(token, journalID, entryID)
	if currentEntryErr != nil {
		return Entry{}, fmt.Errorf("Error obtaining entry (journalID: %s, entryID: %s):\n%s", journalID, entryID, currentEntryErr.Error())
	}

	existingTags := make(map[string]bool)
	for _, tag := range currentEntry.Tags {
		existingTags[tag] = true
	}

	tagsToAdd := make([]string, len(tags))
	numTags := 0
	for _, tag := range tags {
		if _, exists := existingTags[tag]; !exists {
			tagsToAdd[numTags] = tag
			numTags++
		}
	}

	entryTagsRoute := fmt.Sprintf("%s/%s/entries/%s/tags", client.Routes.Journals, journalID, entryID)
	requestBody := entryAddTagsRequest{
		Tags: tagsToAdd[:numTags],
	}

	requestBuffer := new(bytes.Buffer)
	encodeErr := json.NewEncoder(requestBuffer).Encode(requestBody)
	if encodeErr != nil {
		return Entry{}, encodeErr
	}
	request, requestErr := http.NewRequest("POST", entryTagsRoute, requestBuffer)
	if requestErr != nil {
		return Entry{}, requestErr
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	response, responseErr := client.HTTPClient.Do(request)
	if responseErr != nil {
		return Entry{}, responseErr
	}
	defer response.Body.Close()

	// Now return the freshest state of the entry
	return client.GetEntry(token, journalID, entryID)
}

func (client SpireClient) UntagEntry(token, journalID, entryID string, tags []string) (Entry, error) {
	entryTagsRoute := fmt.Sprintf("%s/%s/entries/%s/tags", client.Routes.Journals, journalID, entryID)
	for _, tag := range tags {
		requestBody := entryRemoveTagRequest{Tag: tag}

		requestBuffer := new(bytes.Buffer)
		encodeErr := json.NewEncoder(requestBuffer).Encode(requestBody)
		if encodeErr != nil {
			return Entry{}, encodeErr
		}
		request, requestErr := http.NewRequest("DELETE", entryTagsRoute, requestBuffer)
		if requestErr != nil {
			return Entry{}, requestErr
		}
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Accept", "application/json")
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		response, responseErr := client.HTTPClient.Do(request)
		if responseErr != nil {
			return Entry{}, responseErr
		}
		defer response.Body.Close()
	}

	// Now return the freshest state of the entry
	return client.GetEntry(token, journalID, entryID)
}

func (client SpireClient) UpdateEntry(token, journalID, entryID, title, content string) (Entry, error) {
	// If title or content are empty, does not overwrite the current title and content for the
	// entry with the given entryID. This is different from the behavior of out endpoint.
	currentEntry, currentEntryErr := client.GetEntry(token, journalID, entryID)
	if currentEntryErr != nil {
		return Entry{}, currentEntryErr
	}

	entryRoute := fmt.Sprintf("%s/%s/entries/%s", client.Routes.Journals, journalID, entryID)
	requestBody := entryUpdateRequest{
		Title:   currentEntry.Title,
		Content: currentEntry.Content,
	}
	if title != "" {
		requestBody.Title = title
	}
	if content != "" {
		requestBody.Content = content
	}

	requestBuffer := new(bytes.Buffer)
	encodeErr := json.NewEncoder(requestBuffer).Encode(requestBody)
	if encodeErr != nil {
		return Entry{}, encodeErr
	}
	request, requestErr := http.NewRequest("PUT", entryRoute, requestBuffer)
	if requestErr != nil {
		return Entry{}, requestErr
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	response, responseErr := client.HTTPClient.Do(request)
	if responseErr != nil {
		return Entry{}, responseErr
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return Entry{}, statusErr
	}

	// One more GetEntry API call to get the freshest version of the entry on the server.
	return client.GetEntry(token, journalID, entryID)
}
