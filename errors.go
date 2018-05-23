package buffalo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/x/httpx"
	"github.com/pkg/errors"
)

// HTTPError a typed error returned by http Handlers and used for choosing error handlers
type HTTPError struct {
	Status int   `json:"status"`
	Cause  error `json:"error"`
}

func (h HTTPError) Error() string {
	return h.Cause.Error()
}

// ErrorHandler interface for handling an error for a
// specific status code.
type ErrorHandler func(int, error, Context) error

// ErrorHandlers is used to hold a list of ErrorHandler
// types that can be used to handle specific status codes.
/*
	a.ErrorHandlers[500] = func(status int, err error, c buffalo.Context) error {
		res := c.Response()
		res.WriteHeader(status)
		res.Write([]byte(err.Error()))
		return nil
	}
*/
type ErrorHandlers map[int]ErrorHandler

// Get a registered ErrorHandler for this status code. If
// no ErrorHandler has been registered, a default one will
// be returned.
func (e ErrorHandlers) Get(status int) ErrorHandler {
	if eh, ok := e[status]; ok {
		return eh
	}
	return defaultErrorHandler
}

// PanicHandler recovers from panics gracefully and calls
// the error handling code for a 500 error.
func (a *App) PanicHandler(next Handler) Handler {
	return func(c Context) error {
		defer func() { //catch or finally
			r := recover()
			var err error
			if r != nil { //catch
				switch t := r.(type) {
				case error:
					err = errors.WithStack(t)
				case string:
					err = errors.WithStack(errors.New(t))
				default:
					err = errors.New(fmt.Sprint(t))
				}
				eh := a.ErrorHandlers.Get(500)
				eh(500, err, c)
			}
		}()
		return next(c)
	}
}

func productionErrorResponseFor(status int) []byte {
	if status == http.StatusNotFound {
		return []byte(prodNotFoundTmpl)
	}

	return []byte(prodErrorTmpl)
}

func defaultErrorHandler(status int, origErr error, c Context) error {
	env := c.Value("env")

	c.Logger().Error(origErr)
	c.Response().WriteHeader(status)

	if env != nil && env.(string) == "production" {
		responseBody := productionErrorResponseFor(status)
		c.Response().Write(responseBody)
		return nil
	}

	msg := fmt.Sprintf("%+v", origErr)
	ct := httpx.ContentType(c.Request())
	switch strings.ToLower(ct) {
	case "application/json", "text/json", "json":
		err := json.NewEncoder(c.Response()).Encode(map[string]interface{}{
			"error": msg,
			"code":  status,
		})
		if err != nil {
			return errors.WithStack(err)
		}
	case "application/xml", "text/xml", "xml":
	default:
		if err := c.Request().ParseForm(); err != nil {
			msg = fmt.Sprintf("%s\n%s", err.Error(), msg)
		}
		routes := c.Value("routes")
		if cd, ok := c.(*DefaultContext); ok {
			delete(cd.data, "app")
			delete(cd.data, "routes")
		}
		data := map[string]interface{}{
			"routes":      routes,
			"error":       msg,
			"status":      status,
			"data":        c.Data(),
			"params":      c.Params(),
			"posted_form": c.Request().Form,
			"context":     c,
			"headers":     inspectHeaders(c.Request().Header),
			"inspect": func(v interface{}) string {
				return fmt.Sprintf("%+v", v)
			},
		}
		ctx := plush.NewContextWith(data)
		t, err := plush.Render(devErrorTmpl, ctx)
		if err != nil {
			return errors.WithStack(err)
		}
		res := c.Response()
		_, err = res.Write([]byte(t))
		return err
	}
	return nil
}

type inspectHeaders http.Header

func (i inspectHeaders) String() string {

	bb := make([]string, 0, len(i))

	for k, v := range i {
		bb = append(bb, fmt.Sprintf("%s: %s", k, v))
	}
	sort.Strings(bb)
	return strings.Join(bb, "\n\n")
}
