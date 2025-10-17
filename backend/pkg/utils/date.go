package utils

import (
	"fmt"
	"strings"
	"time"
)

// ParseDate parses a date string in various formats and returns a time.Time
// Supports formats: "MM/DD/YYYY", "MM-DD-YYYY", "YYYY-MM-DD", "YYYY/MM/DD"
func ParseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("date string is empty")
	}

	// List of supported date formats
	formats := []string{
		"02/01/2006", // DD/MM/YYYY
		"02/01/06",   // DD/MM/YY
		"02-01-2006", // DD-MM-YYYY
		"02-01-06",   // DD-MM-YY
		"2/1/2006",   // D/M/YYYY
		"2/1/06",     // D/M/YY
		"2-1-2006",   // D-M-YYYY
		"2-1-06",     // D-M-YY
		"2006-01-02", // YYYY-MM-DD (ISO format)
		"2006/01/02", // YYYY/MM/DD
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date '%s'. Supported formats: DD/MM/YYYY, DD/MM/YY, DD-MM-YYYY, DD-MM-YY, YYYY-MM-DD", dateStr)
}

// ParseDateToPtr parses a date string and returns a pointer to time.Time
func ParseDateToPtr(dateStr string) (*time.Time, error) {
	if strings.TrimSpace(dateStr) == "" {
		return nil, nil
	}

	t, err := ParseDate(dateStr)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
