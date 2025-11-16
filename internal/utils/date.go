package utils

import (
	"fmt"
	"time"
)

// Supported date formats
const (
	// DateFormatYYYYMMDD is the YYYY-MM-DD format
	DateFormatYYYYMMDD = "2006-01-02"
	// DateFormatISO8601 is the ISO 8601 format with time
	DateFormatISO8601 = "2006-01-02T15:04:05Z07:00"
	// DateFormatISO8601Short is the ISO 8601 format without timezone
	DateFormatISO8601Short = "2006-01-02T15:04:05"
)

// ParseDate parses a date string in multiple formats
// Supports:
// - YYYY-MM-DD (e.g., "2024-01-15")
// - ISO 8601 with timezone (e.g., "2024-01-15T10:30:00Z" or "2024-01-15T10:30:00+08:00")
// - ISO 8601 without timezone (e.g., "2024-01-15T10:30:00")
func ParseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("date string is empty")
	}

	// Try parsing in order of most common to least common format
	formats := []string{
		DateFormatYYYYMMDD,     // YYYY-MM-DD
		DateFormatISO8601,      // ISO 8601 with timezone
		DateFormatISO8601Short, // ISO 8601 without timezone
		time.RFC3339,           // RFC3339 (another ISO 8601 variant)
		time.RFC3339Nano,       // RFC3339 with nanoseconds
	}

	var lastErr error
	for _, format := range formats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}

	return time.Time{}, fmt.Errorf("unable to parse date '%s': %w (expected formats: YYYY-MM-DD or ISO 8601)", dateStr, lastErr)
}

// ParseDateToStartOfDay parses a date string and returns the start of that day (00:00:00)
// This is useful for date range queries where you want to include the entire day
func ParseDateToStartOfDay(dateStr string) (time.Time, error) {
	t, err := ParseDate(dateStr)
	if err != nil {
		return time.Time{}, err
	}

	// Truncate to start of day
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()), nil
}

// ParseDateToEndOfDay parses a date string and returns the end of that day (23:59:59)
// This is useful for date range queries where you want to include the entire day
func ParseDateToEndOfDay(dateStr string) (time.Time, error) {
	t, err := ParseDate(dateStr)
	if err != nil {
		return time.Time{}, err
	}

	// Set to end of day
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location()), nil
}

// FormatDate formats a time.Time to YYYY-MM-DD format
func FormatDate(t time.Time) string {
	return t.Format(DateFormatYYYYMMDD)
}

// FormatDateISO8601 formats a time.Time to ISO 8601 format
func FormatDateISO8601(t time.Time) string {
	return t.Format(time.RFC3339)
}
