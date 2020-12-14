package spire

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

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
