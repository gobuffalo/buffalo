package nulls

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

// Float32 adds an implementation for float32
// that supports proper JSON encoding/decoding.
type Float32 struct {
	Float32 float32
	Valid   bool // Valid is true if Float32 is not NULL
}

func (ns Float32) Interface() interface{} {
	if !ns.Valid {
		return nil
	}
	return ns.Float32
}

// NewFloat32 returns a new, properly instantiated
// Float32 object.
func NewFloat32(i float32) Float32 {
	return Float32{Float32: i, Valid: true}
}

// Scan implements the Scanner interface.
func (ns *Float32) Scan(value interface{}) error {
	n := sql.NullFloat64{Float64: float64(ns.Float32)}
	err := n.Scan(value)
	ns.Float32, ns.Valid = float32(n.Float64), n.Valid
	return err
}

// Value implements the driver Valuer interface.
func (ns Float32) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return float64(ns.Float32), nil
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns Float32) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Float32)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the propert representation of that value.
func (ns *Float32) UnmarshalJSON(text []byte) error {
	txt := string(text)
	ns.Valid = true
	if txt == "null" {
		ns.Valid = false
		return nil
	}
	i, err := strconv.ParseFloat(txt, 32)
	if err != nil {
		ns.Valid = false
		return err
	}
	j := float32(i)
	ns.Float32 = j
	return nil
}

func (ns *Float32) UnmarshalText(text []byte) error {
	return ns.UnmarshalJSON(text)
}
