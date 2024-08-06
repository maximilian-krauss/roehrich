package cmd

import (
	"fmt"
	"github.com/maximilian-krauss/roerich/input"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "roehrich",
	Short:             "Tut das not?",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	Args:              cobra.MatchAll(cobra.ExactArgs(1), input.ValidateUrl, cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		println("Hi")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
