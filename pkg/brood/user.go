package brood

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/bugout-dev/bugout-go/pkg/utils"
)

func (client BroodClient) CreateUser(username, email, password string) (User, error) {
	userRoute := client.Routes.User
	data := url.Values{}
	data.Add("username", username)
	data.Add("email", email)
	data.Add("password", password)
	response, err := client.HTTPClient.PostForm(userRoute, data)
	if err != nil {
		return User{}, err
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return User{}, statusErr
	}

	var user User
	decodeErr := json.NewDecoder(response.Body).Decode(&user)
	if decodeErr != nil {
		return user, decodeErr
	}

	return user, nil
}

func (client BroodClient) GenerateToken(username, password string) (string, error) {
	tokenRoute := client.Routes.GenerateToken
	data := url.Values{}
	data.Add("username", username)
	data.Add("password", password)

	response, err := client.HTTPClient.PostForm(tokenRoute, data)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return "", statusErr
	}

	var token UserGeneratedToken
	decodeErr := json.NewDecoder(response.Body).Decode(&token)
	if decodeErr != nil {
		return token.Id, decodeErr
	}

	return token.Id, nil
}

func (client BroodClient) AnnotateToken(token, tokenType, note string) (string, error) {
	tokenRoute := client.Routes.GenerateToken
	data := url.Values{}
	data.Add("access_token", token)
	data.Add("token_type", tokenType)
	data.Add("token_note", note)
	encodedData := data.Encode()

	request, err := http.NewRequest("PUT", tokenRoute, strings.NewReader(encodedData))
	if err != nil {
		return "", err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return "", statusErr
	}

	tokenBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(tokenBytes), nil
}

func (client BroodClient) ListTokens(token string) (UserTokensList, error) {
	listTokensRoute := client.Routes.ListTokens
	request, requestErr := http.NewRequest("GET", listTokensRoute, nil)
	if requestErr != nil {
		return UserTokensList{}, requestErr
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("Accept", "application/json")

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return UserTokensList{}, err
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return UserTokensList{}, err
	}

	var result UserTokensList
	decodeErr := json.NewDecoder(response.Body).Decode(&result)
	return result, decodeErr
}

func (client BroodClient) GetUser(token string) (User, error) {
	userRoute := client.Routes.User
	request, requestErr := http.NewRequest("GET", userRoute, nil)
	if requestErr != nil {
		return User{}, requestErr
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("Accept", "application/json")

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return User{}, err
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return User{}, err
	}

	var user User
	decodeErr := json.NewDecoder(response.Body).Decode(&user)
	return user, decodeErr
}
