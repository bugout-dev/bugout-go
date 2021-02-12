package cmdutils

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func MergeString(stringVar, envName string, onNotSet error) (string, error) {
	finalVal := stringVar

	envVar := os.Getenv(envName)
	if stringVar == "" && envVar != "" {
		finalVal = envVar
	}

	if finalVal == "" && onNotSet != nil {
		return finalVal, onNotSet
	}

	return finalVal, nil
}

const EnvKeyBugoutAccessToken string = "BUGOUT_ACCESS_TOKEN"
const EnvKeyBugoutJournalID string = "BUGOUT_JOURNAL_ID"

var EnvVars []string = []string{
	EnvKeyBugoutAccessToken,
	EnvKeyBugoutJournalID,
}

func IsValidEnvVar(key string) bool {
	for _, validKey := range EnvVars {
		if key == validKey {
			return true
		}
	}

	return false
}

func GenerateArgPopulator(flagName, envName string, required bool) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		flagToken, flagTokenErr := cmd.Flags().GetString(flagName)
		if flagTokenErr != nil {
			return flagTokenErr
		}

		var onNotSet error = nil
		if required {
			onNotSet = fmt.Errorf("Please set the --%s flag or the %s environment variable", flagName, envName)
		}

		finalToken, err := MergeString(flagToken, envName, onNotSet)
		if err != nil {
			return err
		}
		return cmd.Flags().Set(flagName, finalToken)
	}
}

var TokenArgPopulator cobra.PositionalArgs = GenerateArgPopulator("token", EnvKeyBugoutAccessToken, true)
var JournalIDArgPopulator cobra.PositionalArgs = GenerateArgPopulator("journal", EnvKeyBugoutJournalID, true)

func CompositePopulator(populators ...cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		for _, populator := range populators {
			err := populator(cmd, args)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
