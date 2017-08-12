package slices

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// For reading in arrays from postgres

type String []string

func (s String) Interface() interface{} {
	return []string(s)
}

func (s *String) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	(*s) = strToString(string(b))
	return nil
}

func (s String) Value() (driver.Value, error) {
	return fmt.Sprintf("{%s}", strings.Join(s, ",")), nil
}

func (s *String) UnmarshalText(text []byte) error {
	ss := []string{}
	for _, x := range strings.Split(string(text), ",") {
		ss = append(ss, strings.TrimSpace(x))
	}
	(*s) = ss
	return nil
}

func strToString(s string) []string {
	r := strings.Trim(s, "{}")
	return strings.Split(r, ",")
}
