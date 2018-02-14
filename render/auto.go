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

type ErrRedirect struct {
	URL string
}

func (ErrRedirect) Error() string {
	return ""
}

func Auto(ctx context.Context, i interface{}) Renderer {
	e := New(Options{})
	return e.Auto(ctx, i)
}

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
	case "PUT":
		if err := ir.redirect(pname, w, data); err != nil {
			if _, ok := err.(ErrRedirect); ok {
				return err
			}
			return ir.HTML(fmt.Sprintf("%s/edit.html", pname.File())).Render(w, data)
		}
		return nil
	case "POST":
		if err := ir.redirect(pname, w, data); err != nil {
			if _, ok := err.(ErrRedirect); ok {
				return err
			}
			return ir.HTML(fmt.Sprintf("%s/new.html", pname.File())).Render(w, data)
		}
		return nil
	case "DELETE":
		return ErrRedirect{
			URL: "/" + pname.URL(),
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

		return ErrRedirect{
			URL: url,
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
