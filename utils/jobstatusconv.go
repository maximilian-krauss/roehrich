package utils

import "github.com/maximilian-krauss/roehrich/gitlab"

func JobStatusToEmoji(job gitlab.Job) string {
	statusEmojiMap := map[string]string{
		"created":              "⏳",
		"running":              "▶️",
		"failed":               "❌",
		"success":              "✅",
		"canceled":             "⏸️",
		"skipped":              "⏭️",
		"waiting_for_resource": "🕝",
		"manual":               "⚙️",
	}
	converted := statusEmojiMap[job.Status]
	if converted != "" {
		return converted
	}
	return "❓"
}
