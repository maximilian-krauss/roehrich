package utils

import (
	"github.com/fatih/color"
)

func JobStatusToColor(jobStatus string) *color.Color {
	statusColorMap := map[string]color.Attribute{
		"created":              color.FgYellow,
		"pending":              color.FgYellow,
		"running":              color.FgBlue,
		"failed":               color.FgRed,
		"success":              color.FgGreen,
		"canceled":             color.FgRed,
		"skipped":              color.FgWhite,
		"waiting_for_resource": color.FgYellow,
		"manual":               color.FgYellow,
	}
	colorName := statusColorMap[jobStatus]
	if colorName == *new(color.Attribute) {
		return color.New(color.Reset)
	}
	return color.New(colorName, color.Bold)
}
