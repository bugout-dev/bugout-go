package brood

import (
	"encoding/json"
	"net/url"

	"github.com/bugout-dev/bugout-go/pkg/utils"
)

func (client BroodClient) CreateUser(username string, email string, password string) (User, error) {
	userRoute := client.Routes.User
	data := url.Values{}
	data.Add("username", username)
	data.Add("email", email)
	data.Add("password", password)
	response, err := client.HTTPClient.PostForm(userRoute, data)
	if err != nil {
		return User{}, err
	}

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return User{}, statusErr
	}
	defer response.Body.Close()

	var user User
	decodeErr := json.NewDecoder(response.Body).Decode(&user)
	if decodeErr != nil {
		return user, decodeErr
	}

	return user, nil
}
