package brood

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bugout-dev/bugout-go/pkg/utils"
)

// Default Brood URL
const BugoutBroodURL string = "https://auth.bugout.dev"

type BroodCaller interface {
	Ping() (string, error)
	Version() (string, error)
	CreateUser(string, string, string) (User, error)
	GenerateToken(string, string) (string, error)
	AnnotateToken(token, tokenType, note string) (string, error)
	ListTokens(token string) (UserTokensList, error)
	FindUser(token string, queryParameters map[string]string) (User, error)
	GetUser(token string) (User, error)
	VerifyUser(token, code string) (User, error)
	ChangePassword(token, currentPassword, newPassword string) (User, error)
	CreateGroup(token, name string) (Group, error)
	GetUserGroups(token string) (UserGroupsList, error)
	DeleteGroup(token, groupID string) (Group, error)
	RenameGroup(token, groupID, name string) (Group, error)
	AddUserToGroup(token, groupID, username, role string) (UserGroup, error)
	RemoveUserFromGroup(token, groupID, username string) (UserGroup, error)
	CreateResource(token, applicationId string, resourceData interface{}) (Resource, error)
	UpdateResource(token, resourceId string, update interface{}, dropKeys []string) (Resource, error)
	GetResources(token, applicationId string, queryParameters map[string]string) (Resources, error)
	DeleteResource(token, resourceId string) (Resource, error)
	CreateApplication(token, groupId, name, description string) (Application, error)
	GetApplication(token, applicationId string) (Application, error)
	ListApplications(token, groupId string) (ApplicationsList, error)
	DeleteApplication(token, applicationId string) (Application, error)
}

type BroodRoutes struct {
	Ping                string
	Version             string
	User                string
	FindUser            string
	Groups              string
	Token               string
	RevokeToken         string
	ListTokens          string
	ConfirmRegistration string
	ChangePassword      string
	RequestReset        string
	ConfirmReset        string
	Resources           string
	Applications        string
}

func RoutesFromURL(broodURL string) BroodRoutes {
	cleanURL := strings.TrimRight(broodURL, "/")

	return BroodRoutes{
		Ping:                fmt.Sprintf("%s/ping", cleanURL),
		Version:             fmt.Sprintf("%s/version", cleanURL),
		User:                fmt.Sprintf("%s/user", cleanURL),
		FindUser:            fmt.Sprintf("%s/user/find", cleanURL),
		Groups:              fmt.Sprintf("%s/groups", cleanURL),
		Token:               fmt.Sprintf("%s/token", cleanURL),
		ListTokens:          fmt.Sprintf("%s/tokens", cleanURL),
		ConfirmRegistration: fmt.Sprintf("%s/confirm", cleanURL),
		ChangePassword:      fmt.Sprintf("%s/profile/password", cleanURL),
		RequestReset:        fmt.Sprintf("%s/reset", cleanURL),
		ConfirmReset:        fmt.Sprintf("%s/reset_password", cleanURL),
		Resources:           fmt.Sprintf("%s/resources", cleanURL),
		Applications:        fmt.Sprintf("%s/applications", cleanURL),
	}
}

type BroodClient struct {
	BroodURL   string
	Routes     BroodRoutes
	HTTPClient *http.Client
}

func (client BroodClient) Ping() (string, error) {
	pingURL := client.Routes.Ping
	response, err := client.HTTPClient.Get(pingURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return "", statusErr
	}

	pingBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(pingBytes), nil
}

func (client BroodClient) Version() (string, error) {
	versionURL := client.Routes.Version
	response, err := client.HTTPClient.Get(versionURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	statusErr := utils.HTTPStatusCheck(response)
	if statusErr != nil {
		return "", statusErr
	}

	versionBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(versionBytes), nil
}

func NewClient(broodURL string, timeout time.Duration) BroodClient {
	routes := RoutesFromURL(broodURL)
	httpClient := http.Client{Timeout: timeout}
	return BroodClient{BroodURL: broodURL, Routes: routes, HTTPClient: &httpClient}
}

func ClientFromEnv() (BroodClient, error) {
	broodURL := os.Getenv("BUGOUT_BROOD_URL")
	if broodURL == "" {
		broodURL = BugoutBroodURL
	}

	broodTimeoutSecondsRaw := os.Getenv("BUGOUT_TIMEOUT_SECONDS")
	if broodTimeoutSecondsRaw == "" {
		broodTimeoutSecondsRaw = "3"
	}
	broodTimeoutSeconds, conversionErr := strconv.Atoi(broodTimeoutSecondsRaw)
	if conversionErr != nil {
		return BroodClient{}, fmt.Errorf("Could not parse environment variable as integer: BUGOUT_TIMEOUT_SECONDS=%s", broodTimeoutSecondsRaw)
	}
	broodTimeout := time.Duration(broodTimeoutSeconds) * time.Second

	client := NewClient(broodURL, broodTimeout)

	return client, nil
}
