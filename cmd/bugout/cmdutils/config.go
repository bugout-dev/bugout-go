package cmdutils

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func MergeString(stringVar, configKey, flag string, config *viper.Viper, required bool) (string, error) {
	configVal := config.GetString(configKey)
	finalVal := stringVar
	if stringVar == "" && configVal != "" {
		finalVal = configVal
	}

	if finalVal == "" && required {
		return finalVal, fmt.Errorf("Please set %s to use this command", flag)
	}

	return finalVal, nil
}

func TokenArgPopulator(cmd *cobra.Command, args []string) error {
	flagToken, flagTokenErr := cmd.Flags().GetString("token")
	if flagTokenErr != nil {
		return flagTokenErr
	}
	finalToken, err := MergeString(flagToken, "access_token", "-t/--token", viper.GetViper(), true)
	if err != nil {
		return err
	}
	return cmd.Flags().Set("token", finalToken)
}

func JournalIDArgPopulator(cmd *cobra.Command, args []string) error {
	flagToken, flagTokenErr := cmd.Flags().GetString("journal")
	if flagTokenErr != nil {
		return flagTokenErr
	}
	finalToken, err := MergeString(flagToken, "journal_id", "-j/--journal", viper.GetViper(), true)
	if err != nil {
		return err
	}
	return cmd.Flags().Set("journal", finalToken)

}

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
