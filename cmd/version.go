package cmd

import (
	"fmt"
)

const ApplicationVersion = "0.0.10"

func runVersion() {
	fmt.Printf("Current version: %s\n", ApplicationVersion)
}
