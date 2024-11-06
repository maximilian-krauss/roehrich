package cmd

import (
	"fmt"
)

const ApplicationVersion = "0.0.9"

func runVersion() {
	fmt.Printf("Current version: %s\n", ApplicationVersion)
}
