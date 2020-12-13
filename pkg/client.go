package bugout

import "github.com/bugout-dev/bugout-go/pkg/brood"

type BugoutClient struct {
	Brood brood.BroodCaller
}

func ClientFromEnv() (BugoutClient, error) {
	broodClient, err := brood.ClientFromEnv()
	if err != nil {
		return BugoutClient{}, err
	}

	return BugoutClient{Brood: broodClient}, nil
}
