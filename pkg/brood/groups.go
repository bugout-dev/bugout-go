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

func (client BroodClient) RenameGroup(token, groupID, name string) (Group, error) {
	groupsRoute := client.Routes.Groups
	renameRoute := fmt.Sprintf("%s/%s/name", groupsRoute, groupID)
	data := url.Values{}
	data.Add("group_name", name)
	encodedData := data.Encode()

	request, requestErr := http.NewRequest("POST", renameRoute, strings.NewReader(encodedData))
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

func (client BroodClient) AddUserToGroup(token, groupID, username, role string) (UserGroup, error) {
	groupsRoute := client.Routes.Groups
	addUserRoute := fmt.Sprintf("%s/%s/role", groupsRoute, groupID)
	data := url.Values{}
	data.Add("username", username)
	data.Add("user_type", role)
	encodedData := data.Encode()

	request, requestErr := http.NewRequest("POST", addUserRoute, strings.NewReader(encodedData))
	if requestErr != nil {
		return UserGroup{}, requestErr
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("Accept", "application/json")

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return UserGroup{}, err
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return UserGroup{}, statusErr
	}

	var membership UserGroup
	decodeErr := json.NewDecoder(response.Body).Decode(&membership)
	return membership, decodeErr
}

func (client BroodClient) RemoveUserFromGroup(token, groupID, username string) (UserGroup, error) {
	groupsRoute := client.Routes.Groups
	removeUserRoute := fmt.Sprintf("%s/%s/role", groupsRoute, groupID)
	data := url.Values{}
	data.Add("username", username)
	encodedData := data.Encode()

	request, requestErr := http.NewRequest("DELETE", removeUserRoute, strings.NewReader(encodedData))
	if requestErr != nil {
		return UserGroup{}, requestErr
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("Accept", "application/json")

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return UserGroup{}, err
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return UserGroup{}, statusErr
	}

	var membership UserGroup
	decodeErr := json.NewDecoder(response.Body).Decode(&membership)
	return membership, decodeErr
}
