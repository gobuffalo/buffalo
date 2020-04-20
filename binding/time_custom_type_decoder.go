package binding

import (
	"errors"
	"time"
)

type TimeCustomTypeDecoder struct {
	formats *[]string
}

func (td TimeCustomTypeDecoder) Decode(vals []string) (interface{}, error) {
	return td.parseTime(vals)
}

func (td TimeCustomTypeDecoder) parseTime(vals []string) (time.Time, error) {
	var t time.Time
	var err error

	// don't try to parse empty time values, it will raise an error
	if len(vals) == 0 || vals[0] == "" {
		return t, nil
	}

	if td.formats == nil {
		return t, errors.New("empty time format list")
	}

	for _, layout := range *td.formats {
		t, err = time.Parse(layout, vals[0])
		if err == nil {
			return t, nil
		}
	}

	if err != nil {
		return t, err
	}

	return t, nil
}
