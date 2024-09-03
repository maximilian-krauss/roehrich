package utils

func JobStatusToEmoji(jobStatus string) string {
	statusEmojiMap := map[string]string{
		"created":              "⏳",
		"pending":              "⏳",
		"running":              "▶️",
		"failed":               "❌",
		"success":              "✅",
		"canceled":             "⏸️",
		"skipped":              "⏭️",
		"waiting_for_resource": "🕝",
		"manual":               "⚙️",
	}
	converted := statusEmojiMap[jobStatus]
	if converted != "" {
		return converted
	}
	return "❓"
}
