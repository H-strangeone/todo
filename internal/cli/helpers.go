package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/H-strangeone/todo/internal/model"
)

// parseTime parses natural language time strings
func parseTime(s string) (time.Time, error) {
	s = strings.ToLower(strings.TrimSpace(s))

	// Try RFC3339 first
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}

	// Try common formats
	formats := []string{
		"2006-01-02 15:04",
		"2006-01-02 3:04 PM",
		"Jan 02 15:04",
		"Jan 02 3:04 PM",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			// If no year, use current year
			now := time.Now()
			return time.Date(now.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local), nil
		}
	}

	// Handle relative times
	now := time.Now()

	switch {
	case s == "today":
		return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 0, 0, time.Local), nil
	case s == "tomorrow":
		tomorrow := now.Add(24 * time.Hour)
		return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 23, 59, 0, 0, time.Local), nil
	case strings.HasPrefix(s, "tomorrow "):
		// Parse time part after "tomorrow "
		timeStr := strings.TrimPrefix(s, "tomorrow ")
		t, err := parseTimeOfDay(timeStr)
		if err != nil {
			return time.Time{}, err
		}
		tomorrow := now.Add(24 * time.Hour)
		return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), t.Hour(), t.Minute(), 0, 0, time.Local), nil
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", s)
}

// parseTimeOfDay parses time like "5pm", "17:00", "3:30pm"
func parseTimeOfDay(s string) (time.Time, error) {
	formats := []string{
		"3pm",
		"3:04pm",
		"15:04",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time of day: %s", s)
}

// formatReminders formats a slice of durations for display
func formatReminders(reminders model.DurationSlice) string {
	if len(reminders) == 0 {
		return "none"
	}

	strs := make([]string, len(reminders))
	for i, r := range reminders {
		strs[i] = r.String()
	}
	return strings.Join(strs, ", ")
}