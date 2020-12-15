package spire

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	// Have to add trailing slash because of how we set the route on API
	journalsRoute := fmt.Sprintf("%s/", client.Routes.Journals)
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

func ValidMemberTypes() []string {
	return []string{"user", "group"}
}

func IsValidMemberType(memberType string) bool {
	validMemberTypes := ValidMemberTypes()
	for _, validMemberType := range validMemberTypes {
		if memberType == validMemberType {
			return true
		}
	}
	return false
}

func ValidJournalPermissions() []string {
	return []string{"journals.read", "journals.update", "journals.delete"}
}

func IsValidJournalPermission(permission string) bool {
	validJournalPermissions := ValidJournalPermissions()
	for _, validPermission := range validJournalPermissions {
		if permission == validPermission {
			return true
		}
	}

	return false
}

func (client SpireClient) AddJournalMember(token, journalID, memberID, memberType string, permissions []string) (JournalPermissionsList, error) {
	// Validate memberType and permissions
	if !IsValidMemberType(memberType) {
		return JournalPermissionsList{}, fmt.Errorf("Invalid memberType: %s. Choices: %s", memberType, strings.Join(ValidMemberTypes(), ","))
	}

	invalidPermissions := make([]string, len(permissions))
	numInvalidPermissions := 0
	for _, permission := range permissions {
		if !IsValidJournalPermission(permission) {
			invalidPermissions[numInvalidPermissions] = permission
			numInvalidPermissions++
		}
	}
	if numInvalidPermissions > 0 {
		return JournalPermissionsList{}, fmt.Errorf("Invalid permissions: %s. Choices: %s", strings.Join(invalidPermissions[:numInvalidPermissions], ","), strings.Join(ValidJournalPermissions(), ","))
	}

	scopesRoute := fmt.Sprintf("%s/%s/scopes", client.Routes.Journals, journalID)
	requestBody := journalPermissionsRequest{
		HolderID:    memberID,
		HolderType:  memberType,
		Permissions: permissions,
	}
	requestBuffer := new(bytes.Buffer)
	encodeErr := json.NewEncoder(requestBuffer).Encode(requestBody)
	if encodeErr != nil {
		return JournalPermissionsList{}, encodeErr
	}
	request, requestErr := http.NewRequest("POST", scopesRoute, requestBuffer)
	if requestErr != nil {
		return JournalPermissionsList{}, requestErr
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	response, responseErr := client.HTTPClient.Do(request)
	if responseErr != nil {
		return JournalPermissionsList{}, responseErr
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return JournalPermissionsList{}, statusErr
	}

	var permissionsList JournalPermissionsList
	decodeErr := json.NewDecoder(response.Body).Decode(&permissionsList)
	return permissionsList, decodeErr
}

func (client SpireClient) RemoveJournalMember(token, journalID, memberID, memberType string, permissions []string) (JournalPermissionsList, error) {
	// Validate memberType and permissions
	if !IsValidMemberType(memberType) {
		return JournalPermissionsList{}, fmt.Errorf("Invalid memberType: %s. Choices: %s", memberType, strings.Join(ValidMemberTypes(), ","))
	}

	invalidPermissions := make([]string, len(permissions))
	numInvalidPermissions := 0
	for _, permission := range permissions {
		if !IsValidJournalPermission(permission) {
			invalidPermissions[numInvalidPermissions] = permission
			numInvalidPermissions++
		}
	}
	if numInvalidPermissions > 0 {
		return JournalPermissionsList{}, fmt.Errorf("Invalid permissions: %s. Choices: %s", strings.Join(invalidPermissions[:numInvalidPermissions], ","), strings.Join(ValidJournalPermissions(), ","))
	}

	scopesRoute := fmt.Sprintf("%s/%s/scopes", client.Routes.Journals, journalID)
	requestBody := journalPermissionsRequest{
		HolderID:    memberID,
		HolderType:  memberType,
		Permissions: permissions,
	}
	requestBuffer := new(bytes.Buffer)
	encodeErr := json.NewEncoder(requestBuffer).Encode(requestBody)
	if encodeErr != nil {
		return JournalPermissionsList{}, encodeErr
	}
	request, requestErr := http.NewRequest("DELETE", scopesRoute, requestBuffer)
	if requestErr != nil {
		return JournalPermissionsList{}, requestErr
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	response, responseErr := client.HTTPClient.Do(request)
	if responseErr != nil {
		return JournalPermissionsList{}, responseErr
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return JournalPermissionsList{}, statusErr
	}

	var permissionsList JournalPermissionsList
	decodeErr := json.NewDecoder(response.Body).Decode(&permissionsList)
	return permissionsList, decodeErr
}
