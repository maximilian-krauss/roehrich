package cmd

import (
	"github.com/maximilian-krauss/roehrich/input"
	"github.com/maximilian-krauss/roehrich/statuscheck"
	"github.com/spf13/cobra"
	"log"
)

func onlyUrls(_ *cobra.Command, args []string) error {
	maybeUrl := args[0]
	return input.ValidateUrl(maybeUrl)
}

var rootCmd = &cobra.Command{
	Use:               "roehrich",
	Short:             "Tut das not?",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	Args: cobra.MatchAll(
		cobra.ExactArgs(1),
		onlyUrls,
		cobra.OnlyValidArgs,
	),
	Run: func(cmd *cobra.Command, args []string) {
		source := args[0]
		err := statuscheck.Run(statuscheck.Args{
			SourceUrl:                source,
			PollingIntervalInSeconds: 10,
		})
		if err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
