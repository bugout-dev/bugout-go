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

func (client BroodClient) CreateGroup(token, name string) (Group, error) {
	groupsRoute := client.Routes.Groups
	data := url.Values{}
	data.Add("group_name", name)
	encodedData := data.Encode()

	request, requestErr := http.NewRequest("POST", groupsRoute, strings.NewReader(encodedData))
	if requestErr != nil {
		return Group{}, requestErr
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("Accept", "application/json")

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return Group{}, err
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return Group{}, statusErr
	}

	var group Group
	decodeErr := json.NewDecoder(response.Body).Decode(&group)
	return group, decodeErr
}

func (client BroodClient) GetUserGroups(token string) (UserGroupsList, error) {
	groupsRoute := client.Routes.Groups
	request, requestErr := http.NewRequest("GET", groupsRoute, nil)
	if requestErr != nil {
		return UserGroupsList{}, requestErr
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("Accept", "application/json")

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return UserGroupsList{}, err
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return UserGroupsList{}, statusErr
	}

	var userGroups UserGroupsList
	decodeErr := json.NewDecoder(response.Body).Decode(&userGroups)
	return userGroups, decodeErr
}

func (client BroodClient) DeleteGroup(token, groupID string) (Group, error) {
	groupsRoute := client.Routes.Groups
	deletionRoute := fmt.Sprintf("%s/%s", groupsRoute, groupID)

	request, requestErr := http.NewRequest("DELETE", deletionRoute, nil)
	if requestErr != nil {
		return Group{}, requestErr
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("Accept", "application/json")

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return Group{}, err
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return Group{}, statusErr
	}

	var group Group
	decodeErr := json.NewDecoder(response.Body).Decode(&group)
	return group, decodeErr
}
