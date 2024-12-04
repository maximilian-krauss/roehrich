package cmd

import (
	"fmt"
)

const ApplicationVersion = "0.0.10"

func runVersion() error {
	fmt.Printf("Current version: %s\n", ApplicationVersion)
	return nil
}
