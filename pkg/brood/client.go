package brood

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Default Brood URL
const BugoutBroodURL string = "https://auth.bugout.dev"

type BroodCaller interface {
	Ping() (PingResponse, error)
}

type BroodRoutes struct {
	Ping                string
	Version             string
	User                string
	Groups              string
	GenerateToken       string
	RevokeToken         string
	ListTokens          string
	ConfirmRegistration string
	ChangePassword      string
	RequestReset        string
	ConfirmReset        string
}

func RoutesFromURL(broodURL string) BroodRoutes {
	cleanURL := strings.TrimRight(broodURL, "/")

	return BroodRoutes{
		Ping:                fmt.Sprintf("%s/ping", cleanURL),
		Version:             fmt.Sprintf("%s/version", cleanURL),
		User:                fmt.Sprintf("%s/user", cleanURL),
		Groups:              fmt.Sprintf("%s/groups", cleanURL),
		GenerateToken:       fmt.Sprintf("%s/token", cleanURL),
		RevokeToken:         fmt.Sprintf("%s/revoke", cleanURL),
		ListTokens:          fmt.Sprintf("%s/tokens", cleanURL),
		ConfirmRegistration: fmt.Sprintf("%s/confirm", cleanURL),
		ChangePassword:      fmt.Sprintf("%s/profile/password", cleanURL),
		RequestReset:        fmt.Sprintf("%s/reset", cleanURL),
		ConfirmReset:        fmt.Sprintf("%s/reset_password", cleanURL),
	}
}

type BroodClient struct {
	BroodURL   string
	Routes     BroodRoutes
	HTTPClient http.Client
}

func (client *BroodClient) Ping() (PingResponse, error) {
	pingURL := client.Routes.Ping
	response, err := client.HTTPClient.Get(pingURL)
	if err != nil {
		return PingResponse{}, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return PingResponse{}, fmt.Errorf("Invalid status code: %d", response.StatusCode)
	}

	var pingResponse PingResponse
	err = json.NewDecoder(response.Body).Decode(&pingResponse)
	return pingResponse, err
}

func NewClient(broodURL string, timeout time.Duration) BroodClient {
	routes := RoutesFromURL(broodURL)
	httpClient := http.Client{Timeout: timeout}
	return BroodClient{BroodURL: broodURL, Routes: routes, HTTPClient: httpClient}
}

func ClientFromEnv() (BroodClient, error) {
	broodURL := os.Getenv("BUGOUT_BROOD_URL")
	if broodURL == "" {
		broodURL = BugoutBroodURL
	}

	broodTimeoutSecondsRaw := os.Getenv("BUGOUT_BROOD_TIMEOUT_SECONDS")
	if broodTimeoutSecondsRaw == "" {
		broodTimeoutSecondsRaw = "1"
	}
	broodTimeoutSeconds, conversionErr := strconv.Atoi(broodTimeoutSecondsRaw)
	if conversionErr != nil {
		return BroodClient{}, fmt.Errorf("Could not parse environment variable as integer: BUGOUT_BROOD_TIMEOUT_SECONDS=%s", broodTimeoutSecondsRaw)
	}
	broodTimeout := time.Duration(broodTimeoutSeconds) * time.Second

	client := NewClient(broodURL, broodTimeout)

	return client, nil
}
