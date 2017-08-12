// +build !appengine

// Package colortext implements terminal interface for anko script.
package colortext

import (
	"github.com/daviddengcn/go-colortext"
	"github.com/mattn/anko/vm"
)

var ntoc = map[string]ct.Color{
	"none":    ct.None,
	"black":   ct.Black,
	"red":     ct.Red,
	"green":   ct.Green,
	"yellow":  ct.Yellow,
	"blue":    ct.Blue,
	"mazenta": ct.Magenta,
	"cyan":    ct.Cyan,
	"white":   ct.White,
}

func colorOf(name string) ct.Color {
	if c, ok := ntoc[name]; ok {
		return c
	}
	return ct.None
}

func Import(env *vm.Env) *vm.Env {
	m := env.NewPackage("ct")

	m.Define("ChangeColor", func(fg string, fa bool, rest ...interface{}) {
		if len(rest) == 2 {
			bg, ok := rest[0].(string)
			if !ok {
				panic("Argument #3 should be string")
			}
			ba, ok := rest[1].(bool)
			if !ok {
				panic("Argument #4 should be string")
			}
			ct.ChangeColor(colorOf(fg), fa, colorOf(bg), ba)
		} else {
			ct.ChangeColor(colorOf(fg), fa, ct.None, false)
		}
	})

	m.Define("ResetColor", func() {
		ct.ResetColor()
	})
	return m
}
