package slices

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type Map map[string]interface{}

func (m Map) Interface() interface{} {
	return map[string]interface{}(m)
}

func (s *Map) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	err := json.Unmarshal(b, s)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s Map) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return string(b), nil
}

func (m Map) UnmarshalJSON(b []byte) error {
	var stuff map[string]interface{}
	err := json.Unmarshal(b, &stuff)
	if err != nil {
		return err
	}
	for key, value := range stuff {
		m[key] = value
	}
	return nil
}

func (s Map) UnmarshalText(text []byte) error {
	fmt.Println(string(text))
	err := json.Unmarshal(text, &s)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
