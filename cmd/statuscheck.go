package cmd

import (
	"flag"
	"github.com/maximilian-krauss/roehrich/input"
	"github.com/maximilian-krauss/roehrich/statuscheck"
)

func runStatusCheck(args []string) error {
	interval := flag.Int("interval", 10, "Polling interval in seconds")
	flag.Parse()

	source := flag.Arg(0)
	if err := input.ValidateUrl(source); err != nil {
		return err
	}
	err := statuscheck.Run(statuscheck.Args{
		PollingIntervalInSeconds: *interval,
		SourceUrl:                source})

	if err != nil {
		return err
	}
	return nil
}
