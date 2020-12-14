package bugout

import (
	"github.com/bugout-dev/bugout-go/pkg/brood"
	"github.com/bugout-dev/bugout-go/pkg/spire"
)

type BugoutClient struct {
	Brood brood.BroodCaller
	Spire spire.SpireCaller
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
