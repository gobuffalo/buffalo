package nulls

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// Bool replaces sql.NullBool with an implementation
// that supports proper JSON encoding/decoding.
type Bool struct {
	Bool  bool
	Valid bool
}

func (ns Bool) Interface() interface{} {
	if !ns.Valid {
		return nil
	}
	return ns.Bool
}

// NewBool returns a new, properly instantiated
// Boll object.
func NewBool(b bool) Bool {
	return Bool{Bool: b, Valid: true}
}

// Scan implements the Scanner interface.
func (ns *Bool) Scan(value interface{}) error {
	n := sql.NullBool{Bool: ns.Bool}
	err := n.Scan(value)
	ns.Bool, ns.Valid = n.Bool, n.Valid
	return err
}

// Value implements the driver Valuer interface.
func (ns Bool) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Bool, nil
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns Bool) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Bool)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the proper representation of that value. The strings
// "true" and "t" will be considered "true", "false" and "f" will
// be treated as "false". All other values will
//be set to null by Valid = false
func (ns *Bool) UnmarshalJSON(text []byte) error {
	t := string(text)
	if t == "true" || t == "t" {
		ns.Valid = true
		ns.Bool = true
		return nil
	}
	if t == "false" || t == "f" {
		ns.Valid = true
		ns.Bool = false
		return nil
	}
	ns.Bool = false
	ns.Valid = false
	return nil
}

func (ns *Bool) UnmarshalText(text []byte) error {
	return ns.UnmarshalJSON(text)
}
