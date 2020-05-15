package decoders

import (
	"sync"
	"time"
)

var (
	lock = &sync.RWMutex{}

	// timeFormats are the base time formats supported by the time.Time and
	// nulls.Time Decoders you can prepend custom formats to this list
	// by using RegisterTimeFormats.
	timeFormats = []string{
		time.RFC3339,
		"01/02/2006",
		"2006-01-02",
		"2006-01-02T15:04",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}
)

// RegisterTimeFormats allows to add custom time layouts that
// the binder will be able to use for decoding.
func RegisterTimeFormats(layouts ...string) {
	lock.Lock()
	defer lock.Unlock()

	timeFormats = append(layouts, timeFormats...)
}
