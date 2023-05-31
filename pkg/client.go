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

func Client(broodURL, spireURL string, broodTimeout, spireTimeout time.Duration) BugoutClient {
	broodClient := brood.NewClient(broodURL, broodTimeout)
	spireClient := spire.NewClient(spireURL, spireTimeout)

	return BugoutClient{Brood: broodClient, Spire: spireClient}
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
