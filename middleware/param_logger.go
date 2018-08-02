package middleware

import (
	"encoding/json"
	"mime/multipart"
	"net/url"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

//ParameterFilterBlackList is the list of parameter names that will be filtered
//from the application logs (see maskSecrets).
//Important: this list will be used in case insensitive.
var ParameterFilterBlackList = []string{
	"Password",
	"PasswordConfirmation",
	"CreditCard",
	"CVC",
}

var filteredIndicator = []string{"[FILTERED]"}

// ParameterLogger logs form and parameter values to the logger
type parameterLogger struct {
	blacklist []string
}

// ParameterLogger logs form and parameter values to the loggers
func ParameterLogger(next buffalo.Handler) buffalo.Handler {
	pl := parameterLogger{
		blacklist: ParameterFilterBlackList,
	}

	return func(c buffalo.Context) error {
		defer func() {
			req := c.Request()
			if req.Method != "GET" {
				if err := pl.logForm(c); err != nil {
					c.Logger().Error(err)
				}
			}

			b, err := json.Marshal(c.Params())
			if err != nil {
				c.Logger().Error(err)
			}

			c.LogField("params", string(b))
		}()

		return next(c)
	}
}

//Middleware is a buffalo middleware function to connect this parameter filterer with buffalo
func (pl parameterLogger) Middleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		defer func() {
			req := c.Request()
			if req.Method != "GET" {
				if err := pl.logForm(c); err != nil {
					c.Logger().Error(err)
				}
			}

			b, err := json.Marshal(c.Params())
			if err != nil {
				c.Logger().Error(err)
			}

			c.LogField("params", string(b))
		}()

		return next(c)
	}
}

func (pl parameterLogger) logForm(c buffalo.Context) error {
	req := c.Request()
	mp := req.MultipartForm
	if mp != nil {
		return pl.multipartParamLogger(mp, c)
	}

	if err := pl.addFormFieldTo(c, req.Form); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (pl parameterLogger) multipartParamLogger(mp *multipart.Form, c buffalo.Context) error {
	uv := url.Values{}
	for k, v := range mp.Value {
		for _, vv := range v {
			uv.Add(k, vv)
		}
	}
	for k, v := range mp.File {
		for _, vv := range v {
			uv.Add(k, vv.Filename)
		}
	}

	if err := pl.addFormFieldTo(c, uv); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (pl parameterLogger) addFormFieldTo(c buffalo.Context, form url.Values) error {
	maskedForm := pl.maskSecrets(form)
	b, err := json.Marshal(maskedForm)

	if err != nil {
		return err
	}

	c.LogField("form", string(b))
	return nil
}

//maskSecrets matches ParameterFilterBlackList against parameters passed in the
//request, and returns a copy of the request parameters replacing blacklisted params
//with [FILTERED].
func (pl parameterLogger) maskSecrets(form url.Values) url.Values {
	if len(pl.blacklist) == 0 {
		pl.blacklist = ParameterFilterBlackList
	}

	copy := url.Values{}
	for key, values := range form {
	blcheck:
		for _, blacklisted := range pl.blacklist {
			copy[key] = values
			if strings.ToUpper(key) == strings.ToUpper(blacklisted) {
				copy[key] = filteredIndicator
				break blcheck
			}

		}
	}
	return copy
}
