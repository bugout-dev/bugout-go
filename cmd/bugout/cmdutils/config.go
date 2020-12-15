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
