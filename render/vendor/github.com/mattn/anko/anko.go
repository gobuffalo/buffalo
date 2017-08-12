// +build !appengine

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/daviddengcn/go-colortext"
	"github.com/mattn/anko/parser"
	"github.com/mattn/anko/vm"
	"github.com/mattn/go-isatty"

	anko_core "github.com/mattn/anko/builtins"
)

const version = "0.0.1"

var (
	fs   = flag.NewFlagSet(os.Args[0], 1)
	line = fs.String("e", "", "One line of program")
	v    = fs.Bool("v", false, "Display version")

	istty = isatty.IsTerminal(os.Stdout.Fd())
)

func colortext(color ct.Color, bright bool, f func()) {
	if istty {
		ct.ChangeColor(color, bright, ct.None, false)
		f()
		ct.ResetColor()
	} else {
		f()
	}
}

func main() {
	fs.Parse(os.Args[1:])
	if *v {
		fmt.Println(version)
		os.Exit(0)
	}

	var (
		code      string
		b         []byte
		reader    *bufio.Reader
		following bool
		source    string
	)

	env := vm.NewEnv()
	interactive := fs.NArg() == 0 && *line == ""

	env.Define("args", fs.Args())

	if interactive {
		reader = bufio.NewReader(os.Stdin)
		source = "typein"
		os.Args = append([]string{os.Args[0]}, fs.Args()...)
	} else {
		if *line != "" {
			b = []byte(*line)
			source = "argument"
		} else {
			var err error
			b, err = ioutil.ReadFile(fs.Arg(0))
			if err != nil {
				colortext(ct.Red, false, func() {
					fmt.Fprintln(os.Stderr, err)
				})
				os.Exit(1)
			}
			env.Define("args", fs.Args()[1:])
			source = filepath.Clean(fs.Arg(0))
		}
		os.Args = fs.Args()
	}

	anko_core.LoadAllBuiltins(env)

	for {
		if interactive {
			colortext(ct.Green, true, func() {
				if following {
					fmt.Print("  ")
				} else {
					fmt.Print("> ")
				}
			})
			var err error
			b, _, err = reader.ReadLine()
			if err != nil {
				break
			}
			if len(b) == 0 {
				continue
			}
			if code != "" {
				code += "\n"
			}
			code += string(b)
		} else {
			code = string(b)
		}

		parser.EnableErrorVerbose()

		stmts, err := parser.ParseSrc(code)

		if interactive {
			if e, ok := err.(*parser.Error); ok {
				es := e.Error()
				if strings.HasPrefix(es, "syntax error: unexpected") {
					if strings.HasPrefix(es, "syntax error: unexpected $end,") {
						following = true
						continue
					}
				} else {
					if e.Pos.Column == len(b) && !e.Fatal {
						println(e.Error())
						following = true
						continue
					}
					if e.Error() == "unexpected EOF" {
						following = true
						continue
					}
				}
			}
		}

		following = false
		code = ""
		v := vm.NilValue

		if err == nil {
			v, err = vm.Run(stmts, env)
		}
		if err != nil {
			colortext(ct.Red, false, func() {
				if e, ok := err.(*vm.Error); ok {
					fmt.Fprintf(os.Stderr, "%s:%d:%d %s\n", source, e.Pos.Line, e.Pos.Column, err)
				} else if e, ok := err.(*parser.Error); ok {
					if e.Filename != "" {
						source = e.Filename
					}
					fmt.Fprintf(os.Stderr, "%s:%d:%d %s\n", source, e.Pos.Line, e.Pos.Column, err)
				} else {
					fmt.Fprintln(os.Stderr, err)
				}
			})

			if interactive {
				continue
			} else {
				os.Exit(1)
			}
		} else {
			if interactive {
				colortext(ct.Black, true, func() {
					if v == vm.NilValue || !v.IsValid() {
						fmt.Println("nil")
					} else {
						s, ok := v.Interface().(fmt.Stringer)
						if v.Kind() != reflect.String && ok {
							fmt.Println(s)
						} else {
							fmt.Printf("%#v\n", v.Interface())
						}
					}
				})
			} else {
				break
			}
		}
	}
}
