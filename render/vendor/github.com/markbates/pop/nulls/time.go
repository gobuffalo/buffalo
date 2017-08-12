package nulls

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Time replaces sql.NullTime with an implementation
// that supports proper JSON encoding/decoding.
type Time struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

func (ns Time) Interface() interface{} {
	if !ns.Valid {
		return nil
	}
	return ns.Time
}

// NewTime returns a new, properly instantiated
// Time object.
func NewTime(t time.Time) Time {
	return Time{Time: t, Valid: true}
}

// Scan implements the Scanner interface.
func (ns *Time) Scan(value interface{}) error {
	ns.Time, ns.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (ns Time) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Time, nil
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns Time) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Time)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the propert representation of that value.
func (ns *Time) UnmarshalJSON(text []byte) error {
	ns.Valid = false
	txt := string(text)
	if txt == "null" || txt == "" {
		return nil
	}

	t := time.Time{}
	err := t.UnmarshalJSON(text)
	if err == nil {
		ns.Time = t
		ns.Valid = true
	}

	return err
}

func (ns *Time) UnmarshalText(text []byte) error {
	return ns.UnmarshalJSON(text)
}
