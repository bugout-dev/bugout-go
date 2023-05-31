package bugout

import (
	"time"

	"github.com/bugout-dev/bugout-go/pkg/brood"
	"github.com/bugout-dev/bugout-go/pkg/spire"
)

type BugoutClient struct {
	Brood brood.BroodCaller
	Spire spire.SpireCaller
}

func ClientBrood(broodURL string, timeout time.Duration) BugoutClient {
	broodClient := brood.NewClient(broodURL, timeout)
	return BugoutClient{Brood: broodClient}
}

func ClientSpire(spireURL string, timeout time.Duration) BugoutClient {
	spireClient := spire.NewClient(spireURL, timeout)
	return BugoutClient{Spire: spireClient}
}

func ClientFromEnv() (BugoutClient, error) {
	broodClient, err := brood.ClientFromEnv()
	if err != nil {
		return BugoutClient{}, err
	}

	spireClient, err := spire.ClientFromEnv()
	if err != nil {
		return BugoutClient{}, err
	}

	return BugoutClient{Brood: broodClient, Spire: spireClient}, nil
}
