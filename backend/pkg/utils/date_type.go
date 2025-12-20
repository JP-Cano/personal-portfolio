package utils

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Date is a custom type that can scan both time.Time and string values from the database.
// This is necessary because SQLite/Turso stores dates as TEXT and the libsql driver
// returns them as strings instead of time.Time.
//
// Usage:
//   - For required date columns: use `Date`
//   - For nullable date columns: use `*Date` (pointer)
//
// The Scan/Value methods are called automatically by database/sql when reading/writing.
// The MarshalJSON/UnmarshalJSON methods are called automatically during JSON encoding/decoding.
type Date struct {
	time.Time
}

// Common date formats that SQLite/Turso might return
var dateFormats = []string{
	"2006-01-02",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05-07:00",     // With timezone offset (e.g., 2025-12-09 00:00:00+00:00)
	"2006-01-02 15:04:05.000-07:00", // With milliseconds and timezone
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05.000Z",
	"2006-01-02 15:04:05.000",
	time.RFC3339,
	time.RFC3339Nano,
}

// Scan implements sql.Scanner - called automatically when reading from database
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case string:
		return d.parseString(v)
	case []byte:
		return d.parseString(string(v))
	default:
		return fmt.Errorf("cannot scan type %T into Date", value)
	}
}

// parseString parses a date string trying multiple formats
func (d *Date) parseString(s string) error {
	if s == "" {
		d.Time = time.Time{}
		return nil
	}
	for _, format := range dateFormats {
		if parsed, err := time.Parse(format, s); err == nil {
			d.Time = parsed
			return nil
		}
	}
	return fmt.Errorf("cannot parse date: %s", s)
}

// Value implements driver.Valuer - called automatically when writing to database
func (d Date) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Time.Format("2006-01-02"), nil
}

// MarshalJSON implements json.Marshaler - called automatically for JSON encoding
func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return fmt.Appendf(nil, `"%s"`, d.Time.Format("2006-01-02")), nil
}

// UnmarshalJSON implements json.Unmarshaler - called automatically for JSON decoding
func (d *Date) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" || str == `""` {
		d.Time = time.Time{}
		return nil
	}
	// Remove quotes
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	return d.parseString(str)
}
