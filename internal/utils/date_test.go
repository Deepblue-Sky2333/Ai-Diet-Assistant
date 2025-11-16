package utils

import (
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{
			name:      "Valid YYYY-MM-DD format",
			input:     "2024-01-15",
			wantError: false,
		},
		{
			name:      "Valid ISO 8601 with Z timezone",
			input:     "2024-01-15T10:30:00Z",
			wantError: false,
		},
		{
			name:      "Valid ISO 8601 with +08:00 timezone",
			input:     "2024-01-15T10:30:00+08:00",
			wantError: false,
		},
		{
			name:      "Valid ISO 8601 without timezone",
			input:     "2024-01-15T10:30:00",
			wantError: false,
		},
		{
			name:      "Valid RFC3339 format",
			input:     "2024-01-15T10:30:00+00:00",
			wantError: false,
		},
		{
			name:      "Empty string",
			input:     "",
			wantError: true,
		},
		{
			name:      "Invalid format",
			input:     "15-01-2024",
			wantError: true,
		},
		{
			name:      "Invalid date",
			input:     "2024-13-45",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDate(tt.input)
			if tt.wantError {
				if err == nil {
					t.Errorf("ParseDate() expected error but got none for input: %s", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ParseDate() unexpected error: %v for input: %s", err, tt.input)
				}
				if result.IsZero() {
					t.Errorf("ParseDate() returned zero time for valid input: %s", tt.input)
				}
			}
		})
	}
}

func TestParseDateToStartOfDay(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantHour  int
		wantMin   int
		wantSec   int
		wantError bool
	}{
		{
			name:      "YYYY-MM-DD should return start of day",
			input:     "2024-01-15",
			wantHour:  0,
			wantMin:   0,
			wantSec:   0,
			wantError: false,
		},
		{
			name:      "ISO 8601 with time should return start of day",
			input:     "2024-01-15T14:30:45Z",
			wantHour:  0,
			wantMin:   0,
			wantSec:   0,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDateToStartOfDay(tt.input)
			if tt.wantError {
				if err == nil {
					t.Errorf("ParseDateToStartOfDay() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("ParseDateToStartOfDay() unexpected error: %v", err)
				}
				if result.Hour() != tt.wantHour || result.Minute() != tt.wantMin || result.Second() != tt.wantSec {
					t.Errorf("ParseDateToStartOfDay() = %02d:%02d:%02d, want %02d:%02d:%02d",
						result.Hour(), result.Minute(), result.Second(),
						tt.wantHour, tt.wantMin, tt.wantSec)
				}
			}
		})
	}
}

func TestParseDateToEndOfDay(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantHour  int
		wantMin   int
		wantSec   int
		wantError bool
	}{
		{
			name:      "YYYY-MM-DD should return end of day",
			input:     "2024-01-15",
			wantHour:  23,
			wantMin:   59,
			wantSec:   59,
			wantError: false,
		},
		{
			name:      "ISO 8601 with time should return end of day",
			input:     "2024-01-15T10:30:00Z",
			wantHour:  23,
			wantMin:   59,
			wantSec:   59,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDateToEndOfDay(tt.input)
			if tt.wantError {
				if err == nil {
					t.Errorf("ParseDateToEndOfDay() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("ParseDateToEndOfDay() unexpected error: %v", err)
				}
				if result.Hour() != tt.wantHour || result.Minute() != tt.wantMin || result.Second() != tt.wantSec {
					t.Errorf("ParseDateToEndOfDay() = %02d:%02d:%02d, want %02d:%02d:%02d",
						result.Hour(), result.Minute(), result.Second(),
						tt.wantHour, tt.wantMin, tt.wantSec)
				}
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	testTime := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)
	expected := "2024-01-15"
	result := FormatDate(testTime)

	if result != expected {
		t.Errorf("FormatDate() = %s, want %s", result, expected)
	}
}

func TestFormatDateISO8601(t *testing.T) {
	testTime := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)
	result := FormatDateISO8601(testTime)

	// Should be in RFC3339 format
	if result != "2024-01-15T10:30:45Z" {
		t.Errorf("FormatDateISO8601() = %s, want 2024-01-15T10:30:45Z", result)
	}
}

func TestParseDateConsistency(t *testing.T) {
	// Test that different formats for the same date parse to the same day
	inputs := []string{
		"2024-01-15",
		"2024-01-15T00:00:00Z",
		"2024-01-15T10:30:00Z",
		"2024-01-15T23:59:59Z",
	}

	var dates []time.Time
	for _, input := range inputs {
		date, err := ParseDateToStartOfDay(input)
		if err != nil {
			t.Fatalf("ParseDateToStartOfDay() failed for %s: %v", input, err)
		}
		dates = append(dates, date)
	}

	// All should be the same day at 00:00:00
	for i := 1; i < len(dates); i++ {
		if !dates[0].Equal(dates[i]) {
			t.Errorf("ParseDateToStartOfDay() inconsistent: %s != %s",
				dates[0].Format(time.RFC3339), dates[i].Format(time.RFC3339))
		}
	}
}
