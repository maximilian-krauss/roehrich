package cmd

import (
	"flag"
	"log"

	"github.com/maximilian-krauss/roehrich/config"
	"github.com/maximilian-krauss/roehrich/input"
	"github.com/maximilian-krauss/roehrich/retry"
)

func runRetry() error {
	configPath := flag.String("config", config.GetDefaultConfigPath(), "Path to roehrich.yaml")
	flag.Parse()

	source := flag.Arg(1)
	log.Println(source)
	if err := input.ValidateUrl(source); err != nil {
		return err
	}

	err := retry.Run(retry.Args{
		SourceUrl:  source,
		ConfigPath: *configPath,
	})

	return err
}
