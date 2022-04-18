package brood

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bugout-dev/bugout-go/pkg/utils"
)

func (client BroodClient) CreateResource(token, applicationId string, resourceData interface{}) (Resource, error) {
	resourcesRoute := client.Routes.Resources
	requestBody := resourceCreateRequest{
		ApplicationId: applicationId,
		ResourceData:  resourceData,
	}
	requestBuffer := new(bytes.Buffer)
	encodeErr := json.NewEncoder(requestBuffer).Encode(requestBody)
	if encodeErr != nil {
		return Resource{}, encodeErr
	}
	request, requestErr := http.NewRequest("POST", resourcesRoute, requestBuffer)
	if requestErr != nil {
		return Resource{}, requestErr
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("Accept", "application/json")

	response, responseErr := client.HTTPClient.Do(request)
	if responseErr != nil {
		return Resource{}, responseErr
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return Resource{}, statusErr
	}

	var resource Resource
	decodeErr := json.NewDecoder(response.Body).Decode(&resource)
	return resource, decodeErr
}

func (client BroodClient) GetResources(token, applicationId string, queryParameters map[string]string) (Resources, error) {
	resourcesRoute := client.Routes.Resources
	request, requestErr := http.NewRequest("GET", resourcesRoute, nil)
	if requestErr != nil {
		return Resources{}, requestErr
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("Accept", "application/json")

	query := request.URL.Query()
	query.Add("application_id", applicationId)
	for k, v := range queryParameters {
		query.Add(k, v)
	}
	request.URL.RawQuery = query.Encode()

	response, responseErr := client.HTTPClient.Do(request)
	if responseErr != nil {
		return Resources{}, responseErr
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return Resources{}, statusErr
	}

	var resources Resources
	decodeErr := json.NewDecoder(response.Body).Decode(&resources)
	return resources, decodeErr
}

func (client BroodClient) DeleteResource(token, resourceId string) (Resource, error) {
	resourcesRoute := fmt.Sprintf("%s/%s", client.Routes.Resources, resourceId)
	request, requestErr := http.NewRequest("DELETE", resourcesRoute, nil)
	if requestErr != nil {
		return Resource{}, requestErr
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("Accept", "application/json")

	response, responseErr := client.HTTPClient.Do(request)
	if responseErr != nil {
		return Resource{}, responseErr
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return Resource{}, statusErr
	}

	var resource Resource
	decodeErr := json.NewDecoder(response.Body).Decode(&resource)

	return resource, decodeErr
}
