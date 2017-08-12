package plush

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/gobuffalo/plush/ast"

	"github.com/markbates/inflect"
	"github.com/pkg/errors"
)

// Helpers contains all of the default helpers for
// These will be available to all templates. You should add
// any custom global helpers to this list.
var Helpers = HelperMap{
	moot: &sync.Mutex{},
}

func init() {
	Helpers.Add("json", toJSONHelper)
	Helpers.Add("jsEscape", template.JSEscapeString)
	Helpers.Add("htmlEscape", htmlEscape)
	Helpers.Add("upcase", strings.ToUpper)
	Helpers.Add("downcase", strings.ToLower)
	Helpers.Add("contentFor", contentForHelper)
	Helpers.Add("contentOf", contentOfHelper)
	Helpers.Add("markdown", markdownHelper)
	Helpers.Add("len", lenHelper)
	Helpers.Add("debug", debugHelper)
	Helpers.Add("inspect", inspectHelper)
	Helpers.Add("range", rangeHelper)
	Helpers.Add("between", betweenHelper)
	Helpers.Add("until", untilHelper)
	Helpers.Add("groupBy", groupByHelper)
	Helpers.Add("form", BootstrapFormHelper)
	Helpers.Add("form_for", BootstrapFormForHelper)
	Helpers.Add("truncate", truncateHelper)
	Helpers.Add("raw", func(s string) template.HTML {
		return template.HTML(s)
	})
	Helpers.AddMany(inflect.Helpers)
}

// HelperContext is an optional last argument to helpers
// that provides the current context of the call, and access
// to an optional "block" of code that can be executed from
// within the helper.
type HelperContext struct {
	*Context
	compiler *compiler
	block    *ast.BlockStatement
}

const helperContextKind = "HelperContext"

// HasBlock returns true if a block is associated with the helper function
func (h HelperContext) HasBlock() bool {
	return h.block != nil
}

// Block executes the block of template associated with
// the helper, think the block inside of an "if" or "each"
// statement.
func (h HelperContext) Block() (string, error) {
	return h.BlockWith(h.Context)
}

// BlockWith executes the block of template associated with
// the helper, think the block inside of an "if" or "each"
// statement, but with it's own context.
func (h HelperContext) BlockWith(ctx *Context) (string, error) {
	octx := h.compiler.ctx
	defer func() { h.compiler.ctx = octx }()
	h.compiler.ctx = ctx

	if h.block == nil {
		return "", errors.New("no block defined")
	}
	i, err := h.compiler.evalBlockStatement(h.block)
	if err != nil {
		return "", err
	}
	bb := &bytes.Buffer{}
	h.compiler.write(bb, i)
	return bb.String(), nil
}

// toJSONHelper converts an interface into a string.
func toJSONHelper(v interface{}) (template.HTML, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return template.HTML(b), nil
}

func lenHelper(v interface{}) int {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return rv.Len()
}

// Debug by verbosely printing out using 'pre' tags.
func debugHelper(v interface{}) template.HTML {
	return template.HTML(fmt.Sprintf("<pre>%+v</pre>", v))
}

func inspectHelper(v interface{}) string {
	return fmt.Sprintf("%+v", v)
}

func envHelper(k string) string {
	return os.Getenv(k)
}

func htmlEscape(s string, help HelperContext) (string, error) {
	var err error
	if help.HasBlock() {
		s, err = help.Block()
	}
	if err != nil {
		return "", err
	}
	return template.HTMLEscapeString(s), nil
}

func truncateHelper(s string, opts map[string]interface{}) string {
	if opts["size"] == nil {
		opts["size"] = 50
	}
	if opts["trail"] == nil {
		opts["trail"] = "..."
	}
	size := opts["size"].(int)
	if len(s) <= size {
		return s
	}
	trail := opts["trail"].(string)
	return s[:size-len(trail)] + trail
}
