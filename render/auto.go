package render

import (
	"context"
	"fmt"
	"io"
	"path"
	"reflect"
	"regexp"
	"strings"

	"errors"

	"github.com/gobuffalo/flect/name"
)

var errNoID = errors.New("no ID on model")

// ErrRedirect indicates to Context#Render that this is a
// redirect and a template shouldn't be rendered.
type ErrRedirect struct {
	Status int
	URL    string
}

func (ErrRedirect) Error() string {
	return ""
}

// Auto figures out how to render the model based information
// about the request and the name of the model. Auto supports
// automatic rendering of HTML, JSON, and XML. Any status code
// give to Context#Render between 300 - 400 will be respected
// by Auto. Other status codes are not.
/*
# Rules for HTML template lookup:
GET /users - users/index.html
GET /users/id - users/show.html
GET /users/new - users/new.html
GET /users/id/edit - users/edit.html
POST /users - (redirect to /users/id or render user/new.html)
PUT /users/edit - (redirect to /users/id or render user/edit.html)
DELETE /users/id - redirect to /users
*/
func Auto(ctx context.Context, i interface{}) Renderer {
	e := New(Options{})
	return e.Auto(ctx, i)
}

// Auto figures out how to render the model based information
// about the request and the name of the model. Auto supports
// automatic rendering of HTML, JSON, and XML. Any status code
// give to Context#Render between 300 - 400 will be respected
// by Auto. Other status codes are not.
/*
# Rules for HTML template lookup:
GET /users - users/index.html
GET /users/id - users/show.html
GET /users/new - users/new.html
GET /users/id/edit - users/edit.html
POST /users - (redirect to /users/id or render user/new.html)
PUT /users/edit - (redirect to /users/id or render user/edit.html)
DELETE /users/id - redirect to /users
*/
func (e *Engine) Auto(ctx context.Context, i interface{}) Renderer {
	ct, _ := ctx.Value("contentType").(string)
	if ct == "" {
		ct = e.DefaultContentType
	}
	ct = strings.TrimSpace(strings.ToLower(ct))

	if strings.Contains(ct, "json") {
		return e.JSON(i)
	}

	if strings.Contains(ct, "xml") {
		return e.XML(i)
	}

	return htmlAutoRenderer{
		Engine: e,
		model:  i,
	}
}

type htmlAutoRenderer struct {
	*Engine
	model interface{}
}

func (htmlAutoRenderer) ContentType() string {
	return "text/html"
}

func (ir htmlAutoRenderer) Render(w io.Writer, data Data) error {
	n := name.New(ir.typeName())
	pname := name.New(n.Pluralize().String())

	if ir.isPlural() {
		data[pname.VarCasePlural().String()] = ir.model
	} else {
		data[n.VarCaseSingle().String()] = ir.model
	}

	templatePrefix := pname.File()
	if pf, ok := data["template_prefix"].(string); ok {
		templatePrefix = name.New(pf)
	}

	switch data["method"] {
	case "PUT", "POST", "DELETE":
		if err := ir.redirect(pname, w, data); err != nil {
			if er, ok := err.(ErrRedirect); ok && er.Status >= 300 && er.Status < 400 {
				return err
			}
			if data["method"] == "PUT" {
				return ir.HTML(fmt.Sprintf("%s/edit.html", templatePrefix)).Render(w, data)
			}
			return ir.HTML(fmt.Sprintf("%s/new.html", templatePrefix)).Render(w, data)
		}
		return nil
	}
	cp, ok := data["current_path"].(string)

	defCase := func() error {
		return ir.HTML(fmt.Sprintf("%s/%s.html", templatePrefix, "index")).Render(w, data)
	}
	if !ok {
		return defCase()
	}

	if strings.HasSuffix(cp, "/edit/") {
		return ir.HTML(fmt.Sprintf("%s/edit.html", templatePrefix)).Render(w, data)
	}
	if strings.HasSuffix(cp, "/new/") {
		return ir.HTML(fmt.Sprintf("%s/new.html", templatePrefix)).Render(w, data)
	}

	x, err := regexp.Compile(fmt.Sprintf("%s/.+", pname.URL()))
	if err != nil {
		return err
	}
	if x.MatchString(cp) {
		return ir.HTML(fmt.Sprintf("%s/show.html", templatePrefix)).Render(w, data)
	}
	return defCase()
}

func (ir htmlAutoRenderer) redirect(name name.Ident, w io.Writer, data Data) error {
	rv := reflect.Indirect(reflect.ValueOf(ir.model))
	f := rv.FieldByName("ID")
	if !f.IsValid() {
		return errNoID
	}

	fi := f.Interface()
	rt := reflect.TypeOf(fi)
	zero := reflect.Zero(rt)
	if fi != zero.Interface() {
		m, ok := data["method"].(string)
		if !ok {
			m = "GET"
		}
		url := fmt.Sprint(data["current_path"])
		id := fmt.Sprint(f.Interface())
		url = strings.TrimSuffix(url, "/")
		switch m {
		case "DELETE":
			url = strings.TrimSuffix(url, id)
		default:
			if !strings.HasSuffix(url, id) {
				url = path.Join(url, id)
			}
		}

		code := 302
		if i, ok := data["status"].(int); ok {
			if i >= 300 {
				code = i
			}
		}
		return ErrRedirect{
			Status: code,
			URL:    url,
		}
	}
	return errNoID
}

func (ir htmlAutoRenderer) typeName() string {
	rv := reflect.Indirect(reflect.ValueOf(ir.model))
	rt := rv.Type()
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		el := rt.Elem()
		return el.Name()
	default:
		return rt.Name()
	}
}

func (ir htmlAutoRenderer) isPlural() bool {
	rv := reflect.Indirect(reflect.ValueOf(ir.model))
	rt := rv.Type()
	switch rt.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		return true
	}
	return false
}
