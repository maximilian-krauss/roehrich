package cmd

import (
	"fmt"

	"github.com/maximilian-krauss/roehrich/static"
	"github.com/maximilian-krauss/roehrich/update"
)

func runVersion() error {
	fmt.Printf("current version: %s\n", static.ApplicationVersion)

	remoteVersion, err := update.FindLatestVersion(static.ApplicationVersion)
	if err != nil {
		fmt.Printf("failed to check for latest release version: %s", err.Error())
	} else if remoteVersion.IsNewer {
		fmt.Printf("latest version: %s (%s)", remoteVersion.Version, remoteVersion.Url)
		fmt.Println()
	}

	return nil
}
