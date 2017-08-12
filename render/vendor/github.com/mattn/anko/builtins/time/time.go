// Package time implements time interface for anko script.
package time

import (
	t "time"

	"github.com/mattn/anko/vm"
)

func Import(env *vm.Env) *vm.Env {
	m := env.NewPackage("time")
	m.Define("After", t.After)
	m.Define("Sleep", t.Sleep)
	m.Define("Tick", t.Tick)
	m.Define("Since", t.Since)
	m.Define("FixedZone", t.FixedZone)
	m.Define("LoadLocation", t.LoadLocation)
	m.Define("NewTicker", t.NewTicker)
	m.Define("Date", t.Date)
	m.Define("Now", t.Now)
	m.Define("Parse", t.Parse)
	m.Define("ParseDuration", t.ParseDuration)
	m.Define("ParseInLocation", t.ParseInLocation)
	m.Define("Unix", t.Unix)
	m.Define("AfterFunc", t.AfterFunc)
	m.Define("NewTimer", t.NewTimer)
	m.Define("Nanosecond", t.Nanosecond)
	m.Define("Microsecond", t.Microsecond)
	m.Define("Millisecond", t.Millisecond)
	m.Define("Second", t.Second)
	m.Define("Minute", t.Minute)
	m.Define("Hour", t.Hour)
	return m
}
