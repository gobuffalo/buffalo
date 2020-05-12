package decoders

import (
	"errors"
	"time"
)

func parseTime(vals []string, formats []string) (time.Time, error) {
	var t time.Time
	var err error

	// don't try to parse empty time values, it will raise an error
	if len(vals) == 0 || vals[0] == "" {
		return t, nil
	}

	if len(formats) == 0 {
		return t, errors.New("empty time format list")
	}

	for _, layout := range formats {
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
