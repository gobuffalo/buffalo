// Package core implements core interface for anko script.
package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/anko/parser"
	"github.com/mattn/anko/vm"

	anko_encoding_json "github.com/mattn/anko/builtins/encoding/json"
	anko_errors "github.com/mattn/anko/builtins/errors"
	anko_flag "github.com/mattn/anko/builtins/flag"
	anko_fmt "github.com/mattn/anko/builtins/fmt"
	anko_io "github.com/mattn/anko/builtins/io"
	anko_io_ioutil "github.com/mattn/anko/builtins/io/ioutil"
	anko_math "github.com/mattn/anko/builtins/math"
	anko_math_big "github.com/mattn/anko/builtins/math/big"
	anko_math_rand "github.com/mattn/anko/builtins/math/rand"
	anko_net "github.com/mattn/anko/builtins/net"
	anko_net_http "github.com/mattn/anko/builtins/net/http"
	anko_net_url "github.com/mattn/anko/builtins/net/url"
	anko_os "github.com/mattn/anko/builtins/os"
	anko_os_exec "github.com/mattn/anko/builtins/os/exec"
	anko_os_signal "github.com/mattn/anko/builtins/os/signal"
	anko_path "github.com/mattn/anko/builtins/path"
	anko_path_filepath "github.com/mattn/anko/builtins/path/filepath"
	anko_regexp "github.com/mattn/anko/builtins/regexp"
	anko_runtime "github.com/mattn/anko/builtins/runtime"
	anko_sort "github.com/mattn/anko/builtins/sort"
	anko_strings "github.com/mattn/anko/builtins/strings"
	anko_time "github.com/mattn/anko/builtins/time"

	anko_colortext "github.com/mattn/anko/builtins/github.com/daviddengcn/go-colortext"
)

// LoadAllBuiltins is a convenience function that loads all defined builtins.
func LoadAllBuiltins(env *vm.Env) {
	Import(env)

	pkgs := map[string]func(env *vm.Env) *vm.Env{
		"encoding/json": anko_encoding_json.Import,
		"errors":        anko_errors.Import,
		"flag":          anko_flag.Import,
		"fmt":           anko_fmt.Import,
		"io":            anko_io.Import,
		"io/ioutil":     anko_io_ioutil.Import,
		"math":          anko_math.Import,
		"math/big":      anko_math_big.Import,
		"math/rand":     anko_math_rand.Import,
		"net":           anko_net.Import,
		"net/http":      anko_net_http.Import,
		"net/url":       anko_net_url.Import,
		"os":            anko_os.Import,
		"os/exec":       anko_os_exec.Import,
		"os/signal":     anko_os_signal.Import,
		"path":          anko_path.Import,
		"path/filepath": anko_path_filepath.Import,
		"regexp":        anko_regexp.Import,
		"runtime":       anko_runtime.Import,
		"sort":          anko_sort.Import,
		"strings":       anko_strings.Import,
		"time":          anko_time.Import,
		"github.com/daviddengcn/go-colortext": anko_colortext.Import,
	}

	env.Define("import", func(s string) interface{} {
		if loader, ok := pkgs[s]; ok {
			m := loader(env)
			return m
		}
		panic(fmt.Sprintf("package '%s' not found", s))
	})
}

// Import defines core language builtins - len, range, println, int64, etc.
func Import(env *vm.Env) *vm.Env {
	env.Define("len", func(v interface{}) int64 {
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Interface {
			rv = rv.Elem()
		}
		if rv.Kind() == reflect.String {
			return int64(len([]byte(rv.String())))
		}
		if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
			panic("Argument #1 should be array")
		}
		return int64(rv.Len())
	})

	env.Define("keys", func(v interface{}) []string {
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Interface {
			rv = rv.Elem()
		}
		if rv.Kind() != reflect.Map {
			panic("Argument #1 should be map")
		}
		keys := []string{}
		mk := rv.MapKeys()
		for _, key := range mk {
			keys = append(keys, key.String())
		}
		return keys
	})

	env.Define("range", func(args ...int64) []int64 {
		if len(args) < 1 {
			panic("Missing arguments")
		}
		if len(args) > 2 {
			panic("Too many arguments")
		}
		var min, max int64
		if len(args) == 1 {
			min = 0
			max = args[0] - 1
		} else {
			min = args[0]
			max = args[1]
		}
		arr := []int64{}
		for i := min; i <= max; i++ {
			arr = append(arr, i)
		}
		return arr
	})

	env.Define("toString", func(v interface{}) string {
		if b, ok := v.([]byte); ok {
			return string(b)
		}
		return fmt.Sprint(v)
	})

	env.Define("toInt", func(v interface{}) int64 {
		nt := reflect.TypeOf(1)
		rv := reflect.ValueOf(v)
		if rv.Type().ConvertibleTo(nt) {
			return rv.Convert(nt).Int()
		}
		if rv.Kind() == reflect.String {
			i, err := strconv.ParseInt(v.(string), 10, 64)
			if err == nil {
				return i
			}
			f, err := strconv.ParseFloat(v.(string), 64)
			if err == nil {
				return int64(f)
			}
		}
		if rv.Kind() == reflect.Bool {
			if v.(bool) {
				return 1
			}
		}
		return 0
	})

	env.Define("toFloat", func(v interface{}) float64 {
		nt := reflect.TypeOf(1.0)
		rv := reflect.ValueOf(v)
		if rv.Type().ConvertibleTo(nt) {
			return rv.Convert(nt).Float()
		}
		if rv.Kind() == reflect.String {
			f, err := strconv.ParseFloat(v.(string), 64)
			if err == nil {
				return f
			}
		}
		if rv.Kind() == reflect.Bool {
			if v.(bool) {
				return 1.0
			}
		}
		return 0.0
	})

	env.Define("toBool", func(v interface{}) bool {
		nt := reflect.TypeOf(true)
		rv := reflect.ValueOf(v)
		if rv.Type().ConvertibleTo(nt) {
			return rv.Convert(nt).Bool()
		}
		if rv.Type().ConvertibleTo(reflect.TypeOf(1.0)) && rv.Convert(reflect.TypeOf(1.0)).Float() > 0.0 {
			return true
		}
		if rv.Kind() == reflect.String {
			s := strings.ToLower(v.(string))
			if s == "y" || s == "yes" {
				return true
			}
			b, err := strconv.ParseBool(s)
			if err == nil {
				return b
			}
		}
		return false
	})

	env.Define("toChar", func(s rune) string {
		return string(s)
	})

	env.Define("toRune", func(s string) rune {
		if len(s) == 0 {
			return 0
		}
		return []rune(s)[0]
	})

	env.Define("toByteSlice", func(s string) []byte {
		return []byte(s)
	})

	env.Define("toRuneSlice", func(s string) []rune {
		return []rune(s)
	})

	env.Define("toBoolSlice", func(v []interface{}) []bool {
		var result []bool
		toSlice(v, &result)
		return result
	})

	env.Define("toFloatSlice", func(v []interface{}) []float64 {
		var result []float64
		toSlice(v, &result)
		return result
	})

	env.Define("toIntSlice", func(v []interface{}) []int64 {
		var result []int64
		toSlice(v, &result)
		return result
	})

	env.Define("toStringSlice", func(v []interface{}) []string {
		var result []string
		toSlice(v, &result)
		return result
	})

	env.Define("toDuration", func(v int64) time.Duration {
		return time.Duration(v)
	})

	env.Define("typeOf", func(v interface{}) string {
		return reflect.TypeOf(v).String()
	})

	env.Define("chanOf", func(t reflect.Type) reflect.Value {
		return reflect.MakeChan(t, 1)
	})

	env.Define("defined", func(s string) bool {
		_, err := env.Get(s)
		return err == nil
	})

	env.Define("load", func(s string) interface{} {
		body, err := ioutil.ReadFile(s)
		if err != nil {
			panic(err)
		}
		scanner := new(parser.Scanner)
		scanner.Init(string(body))
		stmts, err := parser.Parse(scanner)
		if err != nil {
			if pe, ok := err.(*parser.Error); ok {
				pe.Filename = s
				panic(pe)
			}
			panic(err)
		}
		rv, err := vm.Run(stmts, env)
		if err != nil {
			panic(err)
		}
		if rv.IsValid() && rv.CanInterface() {
			return rv.Interface()
		}
		return nil
	})

	env.Define("panic", func(e interface{}) {
		os.Setenv("ANKO_DEBUG", "1")
		panic(e)
	})

	env.Define("print", fmt.Print)
	env.Define("println", fmt.Println)
	env.Define("printf", fmt.Printf)
	env.Define("close", func(e interface{}) {
		reflect.ValueOf(e).Close()
	})

	env.DefineType("int64", int64(0))
	env.DefineType("float64", float64(0.0))
	env.DefineType("bool", true)
	env.DefineType("string", "")
	return env
}

// toSlice takes in a "generic" slice and converts and copies
// it's elements into the typed slice pointed at by ptr.
// Note that this is a costly operation.
func toSlice(from []interface{}, ptr interface{}) {
	// Value of the pointer to the target
	obj := reflect.Indirect(reflect.ValueOf(ptr))
	// We can't just convert from interface{} to whatever the target is (diff memory layout),
	// so we need to create a New slice of the proper type and copy the values individually
	t := reflect.TypeOf(ptr).Elem()
	slice := reflect.MakeSlice(t, len(from), len(from))
	// Copying the data, val is an adressable Pointer of the actual target type
	val := reflect.Indirect(reflect.New(t.Elem()))
	for i := 0; i < len(from); i++ {
		v := reflect.ValueOf(from[i])
		val.Set(v)
		slice.Index(i).Set(v)
	}
	// Ok now assign our slice to the target pointer
	obj.Set(slice)
}
