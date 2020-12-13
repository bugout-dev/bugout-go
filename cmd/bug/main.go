package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bugout-dev/bugout-go/cmd/bug/brood"
)

func CreateBugCommand() *cobra.Command {
	bugCmd := &cobra.Command{
		Use:   "bug",
		Short: "Interact with Bugout from your command line",
		Long: `Bugout: The knowledge base for software teams

The bug utility lets you interact with your Bugout resources from your command line.`,
	}

	broodCmd := brood.CreateBroodCommand()

	bugCmd.AddCommand(broodCmd)

	return bugCmd
}

func main() {
	bugCmd := CreateBugCommand()
	err := bugCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
