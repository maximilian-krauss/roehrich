package cmd

import (
	"errors"
	"log"
)

const USAGE string = `roerich <flags> <gitlab url>

  Available commands:
      version
      retry
  Available flags:
		--interval [seconds]	Specify an interval for polling pipeline updates (default: 10)
		--config PATH	Specify a path to the configuration file (default: $HOME/.roerich.json)
`

func printFatal(err error) {
	log.Print(err)
	log.Fatal(USAGE)
}

func getCmdFn(command string) func() error {
	switch command {
	case "version":
		return runVersion
	case "retry":
		return runRetry
	default:
		return runStatusCheck
	}
}

func Execute(args []string) {
	if len(args) == 0 {
		printFatal(errors.New("no arguments have been provided"))
	}
	command := args[0]
	err := getCmdFn(command)()
	if err != nil {
		printFatal(err)
	}
}
