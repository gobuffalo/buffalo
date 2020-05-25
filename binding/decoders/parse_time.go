package decoders

import (
	"time"
)

func parseTime(vals []string) (time.Time, error) {
	var t time.Time
	var err error

	// don't try to parse empty time values, it will raise an error
	if len(vals) == 0 || vals[0] == "" {
		return t, nil
	}

	for _, layout := range timeFormats {
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
