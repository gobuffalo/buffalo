package render

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"

	"github.com/markbates/inflect"
	"github.com/pkg/errors"
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
	ct, ok := ctx.Value("contentType").(string)
	if !ok {
		ct = "text/html"
	}
	ct = strings.ToLower(ct)

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
	name := inflect.Name(ir.typeName())
	name = inflect.Name(name.Singular())
	pname := inflect.Name(name.Plural())

	if ir.isPlural() {
		data[pname.VarCasePlural()] = ir.model
	} else {
		data[name.VarCaseSingular()] = ir.model
	}

	switch data["method"] {
	case "PUT", "POST":
		if err := ir.redirect(pname, w, data); err != nil {
			if er, ok := err.(ErrRedirect); ok && er.Status >= 300 && er.Status < 400 {
				return err
			}
			if data["method"] == "PUT" {
				return ir.HTML(fmt.Sprintf("%s/edit.html", pname.File())).Render(w, data)
			}
			return ir.HTML(fmt.Sprintf("%s/new.html", pname.File())).Render(w, data)
		}
		return nil
	case "DELETE":
		return ErrRedirect{
			Status: 302,
			URL:    "/" + pname.URL(),
		}
	}
	if cp, ok := data["current_path"].(string); ok {
		if strings.HasSuffix(cp, "/edit") {
			return ir.HTML(fmt.Sprintf("%s/edit.html", pname.File())).Render(w, data)
		}
		if strings.HasSuffix(cp, "/new") {
			return ir.HTML(fmt.Sprintf("%s/new.html", pname.File())).Render(w, data)
		}

		x, err := regexp.Compile(fmt.Sprintf("%s/.+", pname.URL()))
		if err != nil {
			return errors.WithStack(err)
		}
		if x.MatchString(cp) {
			return ir.HTML(fmt.Sprintf("%s/show.html", pname.File())).Render(w, data)
		}
	}

	return ir.HTML(fmt.Sprintf("%s/%s.html", pname.File(), "index")).Render(w, data)
}

func (ir htmlAutoRenderer) redirect(name inflect.Name, w io.Writer, data Data) error {
	rv := reflect.Indirect(reflect.ValueOf(ir.model))
	f := rv.FieldByName("ID")
	if !f.IsValid() {
		return errNoID
	}

	fi := f.Interface()
	rt := reflect.TypeOf(fi)
	zero := reflect.Zero(rt)
	if fi != zero.Interface() {
		url := fmt.Sprintf("/%s/%v", name.URL(), f.Interface())

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
