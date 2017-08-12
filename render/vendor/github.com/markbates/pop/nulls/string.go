package nulls

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// String replaces sql.NullString with an implementation
// that supports proper JSON encoding/decoding.
type String sql.NullString

func (ns String) Interface() interface{} {
	if !ns.Valid {
		return nil
	}
	return ns.String
}

// NewString returns a new, properly instantiated
// String object.
func NewString(s string) String {
	return String{String: s, Valid: true}
}

// Scan implements the Scanner interface.
func (ns *String) Scan(value interface{}) error {
	n := sql.NullString{String: ns.String}
	err := n.Scan(value)
	ns.String, ns.Valid = n.String, n.Valid
	return err
}

// Value implements the driver Valuer interface.
func (ns String) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns String) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the propert representation of that value.
func (ns *String) UnmarshalJSON(text []byte) error {
	ns.Valid = false
	if string(text) == "null" {
		return nil
	}
	if err := json.Unmarshal(text, &ns.String); err == nil {
		ns.Valid = true
	}
	return nil
}

func (ns *String) UnmarshalText(text []byte) error {
	ns.Valid = false
	t := string(text)
	if t == "null" {
		return nil
	}
	ns.String = t
	ns.Valid = true
	return nil
}
