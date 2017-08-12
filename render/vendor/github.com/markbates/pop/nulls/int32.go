package nulls

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

// Int32 adds an implementation for int32
// that supports proper JSON encoding/decoding.
type Int32 struct {
	Int32 int32
	Valid bool // Valid is true if Int32 is not NULL
}

func (ns Int32) Interface() interface{} {
	if !ns.Valid {
		return nil
	}
	return ns.Int32
}

// NewInt32 returns a new, properly instantiated
// Int object.
func NewInt32(i int32) Int32 {
	return Int32{Int32: i, Valid: true}
}

// Scan implements the Scanner interface.
func (ns *Int32) Scan(value interface{}) error {
	n := sql.NullInt64{Int64: int64(ns.Int32)}
	err := n.Scan(value)
	ns.Int32, ns.Valid = int32(n.Int64), n.Valid
	return err
}

// Value implements the driver Valuer interface.
func (ns Int32) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return int64(ns.Int32), nil
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns Int32) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Int32)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the propert representation of that value.
func (ns *Int32) UnmarshalJSON(text []byte) error {
	txt := string(text)
	ns.Valid = true
	if txt == "null" {
		ns.Valid = false
		return nil
	}
	i, err := strconv.ParseInt(txt, 10, 32)
	if err != nil {
		ns.Valid = false
		return err
	}
	j := int32(i)
	ns.Int32 = j
	return nil
}

func (ns *Int32) UnmarshalText(text []byte) error {
	return ns.UnmarshalJSON(text)
}
