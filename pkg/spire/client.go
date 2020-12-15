package spire

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

// Default Spire URL
const BugoutSpireURL string = "https://spire.bugout.dev"

type SpireCaller interface {
	Ping() (string, error)
	CreateJournal(token, name string) (Journal, error)
	GetJournal(token, journalID string) (Journal, error)
	ListJournals(token string) (JournalsList, error)
	UpdateJournal(token, journalID, name string) (Journal, error)
	DeleteJournal(token, journalID string) (Journal, error)
	AddJournalMember(token, journalID, memberID, memberType string, permissions []string) (JournalPermissionsList, error)
	RemoveJournalMember(token, journalID, memberID, memberType string, permissions []string) (JournalPermissionsList, error)
	CreateEntry(token, journalID, title, content string, tags []string, context EntryContext) (Entry, error)
	DeleteEntry(token, journalID, entryID string) (Entry, error)
	GetEntry(token, journalID, entryID string) (Entry, error)
	ListEntries(token, journalID string, limit, offset int) (EntryResultsPage, error)
	SearchEntries(token, journalID, searchQuery string, limit, offset int) (EntryResultsPage, error)
	TagEntry(token, journalID, entryID string, tags []string) (Entry, error)
	UntagEntry(token, journalID, entryID string, tags []string) (Entry, error)
	UpdateEntry(token, journalID, entryID, title, content string) (Entry, error)
}

type SpireRoutes struct {
	Ping     string
	Journals string
}

func RoutesFromURL(spireURL string) SpireRoutes {
	cleanURL := strings.TrimRight(spireURL, "/")

	return SpireRoutes{
		Ping:     fmt.Sprintf("%s/ping", cleanURL),
		Journals: fmt.Sprintf("%s/journals", cleanURL),
	}
}

type SpireClient struct {
	SpireURL   string
	Routes     SpireRoutes
	HTTPClient *http.Client
}

func (client SpireClient) Ping() (string, error) {
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

func NewClient(spireURL string, timeout time.Duration) SpireClient {
	routes := RoutesFromURL(spireURL)
	httpClient := http.Client{Timeout: timeout}
	return SpireClient{SpireURL: spireURL, Routes: routes, HTTPClient: &httpClient}
}

func ClientFromEnv() (SpireClient, error) {
	spireURL := os.Getenv("BUGOUT_SPIRE_URL")
	if spireURL == "" {
		spireURL = BugoutSpireURL
	}

	spireTimeoutSecondsRaw := os.Getenv("BUGOUT_SPIRE_TIMEOUT_SECONDS")
	if spireTimeoutSecondsRaw == "" {
		spireTimeoutSecondsRaw = "1"
	}
	spireTimeoutSeconds, conversionErr := strconv.Atoi(spireTimeoutSecondsRaw)
	if conversionErr != nil {
		return SpireClient{}, fmt.Errorf("Could not parse environment variable as integer: BUGOUT_BROOD_TIMEOUT_SECONDS=%s", spireTimeoutSecondsRaw)
	}
	spireTimeout := time.Duration(spireTimeoutSeconds) * time.Second

	client := NewClient(spireURL, spireTimeout)

	return client, nil
}
