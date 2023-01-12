package brood

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/bugout-dev/bugout-go/pkg/utils"
)

func (client BroodClient) CreateApplication(token, groupId, name, description string) (Application, error) {
	applicationsRoute := client.Routes.Applications
	data := url.Values{}
	data.Add("group_id", groupId)
	data.Add("name", name)
	if description != "" {
		data.Add("description", description)
	}
	encodedData := data.Encode()

	request, requestErr := http.NewRequest("POST", applicationsRoute, strings.NewReader(encodedData))
	if requestErr != nil {
		return Application{}, requestErr
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("Accept", "application/json")
	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return Application{}, err
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return Application{}, statusErr
	}

	var application Application
	decodeErr := json.NewDecoder(response.Body).Decode(&application)
	return application, decodeErr
}

func (client BroodClient) GetApplication(token, applicationId string) (Application, error) {
	applicationsRoute := client.Routes.Applications
	specificApplicationRoute := fmt.Sprintf("%s/%s", applicationsRoute, applicationId)
	request, requestErr := http.NewRequest("GET", specificApplicationRoute, nil)
	if requestErr != nil {
		return Application{}, requestErr
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("Accept", "application/json")

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return Application{}, err
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return Application{}, statusErr
	}

	var application Application
	decodeErr := json.NewDecoder(response.Body).Decode(&application)
	return application, decodeErr
}

func (client BroodClient) ListApplications(token, groupId string) (ApplicationsList, error) {
	applicationsRoute := client.Routes.Applications
	request, requestErr := http.NewRequest("GET", applicationsRoute, nil)
	if requestErr != nil {
		return ApplicationsList{}, requestErr
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("Accept", "application/json")

	query := request.URL.Query()
	if groupId != "" {
		query.Add("group_id", groupId)
	}
	request.URL.RawQuery = query.Encode()

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return ApplicationsList{}, err
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return ApplicationsList{}, statusErr
	}

	var applications ApplicationsList
	decodeErr := json.NewDecoder(response.Body).Decode(&applications)
	return applications, decodeErr
}

func (client BroodClient) DeleteApplication(token, applicationId string) (Application, error) {
	applicationsRoute := client.Routes.Applications
	deletionRoute := fmt.Sprintf("%s/%s", applicationsRoute, applicationId)

	request, requestErr := http.NewRequest("DELETE", deletionRoute, nil)
	if requestErr != nil {
		return Application{}, requestErr
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("Accept", "application/json")

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return Application{}, err
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return Application{}, statusErr
	}

	var application Application
	decodeErr := json.NewDecoder(response.Body).Decode(&application)
	return application, decodeErr
}
