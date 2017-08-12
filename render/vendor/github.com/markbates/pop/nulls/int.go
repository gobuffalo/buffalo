package nulls

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

// Int adds an implementation for int
// that supports proper JSON encoding/decoding.
type Int struct {
	Int   int
	Valid bool // Valid is true if Int is not NULL
}

func (ns Int) Interface() interface{} {
	if !ns.Valid {
		return nil
	}
	return ns.Int
}

// NewInt returns a new, properly instantiated
// Int object.
func NewInt(i int) Int {
	return Int{Int: i, Valid: true}
}

// Scan implements the Scanner interface.
func (ns *Int) Scan(value interface{}) error {
	n := sql.NullInt64{Int64: int64(ns.Int)}
	err := n.Scan(value)
	ns.Int, ns.Valid = int(n.Int64), n.Valid
	return err
}

// Value implements the driver Valuer interface.
func (ns Int) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return int64(ns.Int), nil
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns Int) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Int)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the propert representation of that value.
func (ns *Int) UnmarshalJSON(text []byte) error {
	if i, err := strconv.ParseInt(string(text), 10, strconv.IntSize); err == nil {
		ns.Valid = true
		ns.Int = int(i)
	}
	return nil
}

func (ns *Int) UnmarshalText(text []byte) error {
	return ns.UnmarshalJSON(text)
}
