package cli

import (
	"fmt"
	"time"
)

// formatDuration formats a duration in a human-readable format
func formatDuration(d time.Duration) string {
	if d < 0 {
		return "Completed"
	}
	if d >= time.Hour {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		if minutes == 0 {
			return fmt.Sprintf("%d hour%s", hours, plural(hours))
		}
		return fmt.Sprintf("%d hour%s %d minute%s", hours, plural(hours), minutes, plural(minutes))
	}
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	if seconds == 0 {
		return fmt.Sprintf("%d minute%s", minutes, plural(minutes))
	}
	return fmt.Sprintf("%d minute%s %d second%s", minutes, plural(minutes), seconds, plural(seconds))
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
