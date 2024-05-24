package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:    "röhrich",
		Version: "0.0.1",
		Usage:   "Keeps track of a merge request pipeline",
		Action: func(cliContext *cli.Context) error {
			mrUrl := cliContext.Args().Get(0)
			if mrUrl == "" {
				return fmt.Errorf("please provide a merge request url")
			}
			fmt.Printf("Merge request URL: %s\n", mrUrl)
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
