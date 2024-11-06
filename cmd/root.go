package cmd

import (
	"github.com/pkg/errors"
	"log"
)

const USAGE string = `roerich <flags> <gitlab url>
  Available flags:
	- TODO
`

func printFatal(err error) {
	log.Print(err)
	log.Fatal(USAGE)
}

func Execute(args []string) {
	if len(args) == 0 {
		printFatal(errors.New("no arguments have been provided"))
	}
	command := args[0]

	switch command {
	case "version":
		runVersion()
	default:
		err := runStatusCheck()
		if err != nil {
			printFatal(err)
		}
	}
}
