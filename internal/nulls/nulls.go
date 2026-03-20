package nulls

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Time represents a time.Time that may be null.
// It implements sql.Scanner and driver.Valuer interfaces.
type Time struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the sql.Scanner interface.
func (t *Time) Scan(value interface{}) error {
	if value == nil {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}
	t.Valid = true
	switch v := value.(type) {
	case time.Time:
		t.Time = v
	case []byte:
		return t.Parse(string(v))
	case string:
		return t.Parse(v)
	}
	return nil
}

// Parse tries to parse the string as a time using multiple formats.
func (t *Time) Parse(s string) error {
	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
		"01/02/2006",
		"01/02/2006 15:04:05",
		"2006-01-02T15:04:05",
		time.RFC3339Nano,
	}

	for _, format := range formats {
		if tt, err := time.Parse(format, s); err == nil {
			t.Time = tt
			return nil
		}
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}
	return t.Parse(string(data))
}

// MarshalJSON implements json.Marshaler.
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time)
}

// String implements fmt.Stringer.
func (t Time) String() string {
	if !t.Valid {
		return ""
	}
	return t.Time.String()
}

// NewTime returns a new, properly initialized
// Time object.
func NewTime(t time.Time) Time {
	return Time{
		Time:  t,
		Valid: true,
	}
}
