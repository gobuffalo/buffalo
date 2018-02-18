package render

import (
	"context"
	"fmt"
	"io"
	"path"
	"reflect"
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
	name := inflect.Name(ir.typeName().Singular())
	pname := inflect.Name(name.Plural())

	if ir.isPlural() {
		data[pname.VarCasePlural()] = ir.model
	} else {
		data[name.VarCaseSingular()] = ir.model
	}

	cp, ok := data["current_path"].(string)
	switch data["method"] {
	case "PUT":
		code := ir.status(data)
		// if successful redirect to the GET version of the URL
		// PUT /users/1 -> redirect -> GET /users/1
		if code < 400 {
			return ErrRedirect{
				Status: code,
				URL:    cp,
			}
		}
		if ok {
			// PUT /users/1 -> /users
			cp = path.Dir(cp)
		} else {
			cp = pname.File()
		}
		return ir.HTML(fmt.Sprintf("%s/edit.html", cp)).Render(w, data)
	case "POST":
		if err := ir.redirect(cp, w, data); err != nil {
			if er, ok := err.(ErrRedirect); ok && er.Status >= 300 && er.Status < 400 {
				return err
			}
		}
		return ir.HTML(fmt.Sprintf("%s/new.html", cp)).Render(w, data)
	case "DELETE":
		if ok {
			// DELETE /users/{id} -> /users
			cp = path.Dir(cp)
		} else {
			cp = "/" + pname.URL()
		}
		return ErrRedirect{
			Status: 302,
			URL:    cp,
		}
	}
	if ok {
		if strings.HasSuffix(cp, "/edit") {
			// GET /users/{id}/edit -> /users
			cp = path.Dir(path.Dir(cp))
			return ir.HTML(fmt.Sprintf("%s/edit.html", cp)).Render(w, data)
		}
		if strings.HasSuffix(cp, "/new") {
			// GET /users/new -> /users
			cp = path.Dir(cp)
			return ir.HTML(fmt.Sprintf("%s/new.html", cp)).Render(w, data)
		}

		if ir.isPlural() {
			// GET /users - if it's a slice/array render the index page
			return ir.HTML(fmt.Sprintf("%s/%s.html", cp, "index")).Render(w, data)
		}
		// GET /users/{id}
		return ir.HTML(fmt.Sprintf("%s/show.html", path.Dir(cp))).Render(w, data)
	}

	return errors.New("could not auto render this model, please render it manually")
}

func (ir htmlAutoRenderer) redirect(path string, w io.Writer, data Data) error {
	rv := reflect.Indirect(reflect.ValueOf(ir.model))
	f := rv.FieldByName("ID")
	if !f.IsValid() {
		return errNoID
	}

	fi := f.Interface()
	rt := reflect.TypeOf(fi)
	zero := reflect.Zero(rt)
	if fi != zero.Interface() {
		url := fmt.Sprintf("%s/%v", path, f.Interface())

		return ErrRedirect{
			Status: ir.status(data),
			URL:    url,
		}
	}
	return errNoID
}

func (ir htmlAutoRenderer) status(data Data) int {
	if i, ok := data["status"].(int); ok {
		if i >= 300 {
			return i
		}
	}
	return 302
}

func (ir htmlAutoRenderer) typeName() inflect.Name {
	rv := reflect.Indirect(reflect.ValueOf(ir.model))
	rt := rv.Type()
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		el := rt.Elem()
		return inflect.Name(el.Name())
	}
	return inflect.Name(rt.Name())
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
