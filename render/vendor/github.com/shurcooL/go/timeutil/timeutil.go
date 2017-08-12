// Package timeutil provides a func for getting start of week of given time.
package timeutil

import "time"

// StartOfDay returns time at start of day of t.
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// StartOfWeek returns time at start of week of t.
func StartOfWeek(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day()-int(t.Weekday()), 0, 0, 0, 0, t.Location())
}
