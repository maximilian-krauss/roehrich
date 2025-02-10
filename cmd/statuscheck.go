package cmd

import (
	"flag"
	"github.com/maximilian-krauss/roehrich/config"
	"github.com/maximilian-krauss/roehrich/input"
	"github.com/maximilian-krauss/roehrich/statuscheck"
)

func runStatusCheck() error {
	interval := flag.Int("interval", 10, "Polling interval in seconds")
	configPath := flag.String("config", config.GetDefaultConfigPath(), "Path to roehrich.yaml")
	skipVersionCheck := flag.Bool("skip-version-check", false, "Skip version check")
	flag.Parse()

	source := flag.Arg(0)
	if err := input.ValidateUrl(source); err != nil {
		return err
	}
	err := statuscheck.Run(statuscheck.Args{
		PollingIntervalInSeconds: *interval,
		SourceUrl:                source,
		ConfigPath:               *configPath,
		SkipVersionCheck:         *skipVersionCheck,
	})

	if err != nil {
		return err
	}
	return nil
}
