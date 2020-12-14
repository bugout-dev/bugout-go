package spire

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bugout-dev/bugout-go/pkg/utils"
)

func (client SpireClient) CreateJournal(token, name string) (Journal, error) {
	journalsRoute := client.Routes.Journals
	requestBody := journalCreateRequest{Name: name}
	requestBuffer := new(bytes.Buffer)
	encodeErr := json.NewEncoder(requestBuffer).Encode(requestBody)
	if encodeErr != nil {
		return Journal{}, encodeErr
	}
	request, requestErr := http.NewRequest("POST", journalsRoute, requestBuffer)
	if requestErr != nil {
		return Journal{}, requestErr
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	response, responseErr := client.HTTPClient.Do(request)
	if responseErr != nil {
		return Journal{}, responseErr
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return Journal{}, statusErr
	}

	var journal Journal
	decodeErr := json.NewDecoder(response.Body).Decode(&journal)
	return journal, decodeErr
}

func (client SpireClient) GetJournal(token, journalID string) (Journal, error) {
	journalRoute := fmt.Sprintf("%s/%s", client.Routes.Journals, journalID)
	request, requestErr := http.NewRequest("GET", journalRoute, nil)
	if requestErr != nil {
		return Journal{}, requestErr
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	response, responseErr := client.HTTPClient.Do(request)
	if responseErr != nil {
		return Journal{}, responseErr
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return Journal{}, statusErr
	}

	var journal Journal
	decodeErr := json.NewDecoder(response.Body).Decode(&journal)
	return journal, decodeErr
}

func (client SpireClient) ListJournals(token string) (JournalsList, error) {
	journalsRoute := client.Routes.Journals
	request, requestErr := http.NewRequest("GET", journalsRoute, nil)
	if requestErr != nil {
		return JournalsList{}, requestErr
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	response, responseErr := client.HTTPClient.Do(request)
	if responseErr != nil {
		return JournalsList{}, responseErr
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return JournalsList{}, statusErr
	}

	var journals JournalsList
	decodeErr := json.NewDecoder(response.Body).Decode(&journals)
	return journals, decodeErr
}

func (client SpireClient) UpdateJournal(token, journalID, name string) (Journal, error) {
	journalRoute := fmt.Sprintf("%s/%s", client.Routes.Journals, journalID)
	requestBody := journalCreateRequest{Name: name}
	requestBuffer := new(bytes.Buffer)
	encodeErr := json.NewEncoder(requestBuffer).Encode(requestBody)
	if encodeErr != nil {
		return Journal{}, encodeErr
	}
	request, requestErr := http.NewRequest("PUT", journalRoute, requestBuffer)
	if requestErr != nil {
		return Journal{}, requestErr
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	response, responseErr := client.HTTPClient.Do(request)
	if responseErr != nil {
		return Journal{}, responseErr
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return Journal{}, statusErr
	}

	var journal Journal
	decodeErr := json.NewDecoder(response.Body).Decode(&journal)
	return journal, decodeErr
}

func (client SpireClient) DeleteJournal(token, journalID string) (Journal, error) {
	journalRoute := fmt.Sprintf("%s/%s", client.Routes.Journals, journalID)
	request, requestErr := http.NewRequest("DELETE", journalRoute, nil)
	if requestErr != nil {
		return Journal{}, requestErr
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	response, responseErr := client.HTTPClient.Do(request)
	if responseErr != nil {
		return Journal{}, responseErr
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return Journal{}, statusErr
	}

	var journal Journal
	decodeErr := json.NewDecoder(response.Body).Decode(&journal)
	return journal, decodeErr
}
