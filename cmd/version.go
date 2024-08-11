package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const ApplicationVersion = "0.0.3"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the current application version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Current version: %s\n", ApplicationVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
