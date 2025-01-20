package cmd

import (
	"fmt"

	"github.com/maximilian-krauss/roehrich/update"
)

const ApplicationVersion = "0.0.15"

func runVersion() error {
	fmt.Printf("current version: %s\n", ApplicationVersion)

	remoteVersion, err := update.FindLatestVersion(ApplicationVersion)
	if err != nil {
		fmt.Printf("failed to check for latest release version: %s", err.Error())
	} else if remoteVersion.IsNewer {
		fmt.Printf("latest version: %s (%s)", remoteVersion.Version, remoteVersion.Url)
		fmt.Println()
	}

	return nil
}
