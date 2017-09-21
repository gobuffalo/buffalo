package buffalo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gobuffalo/plush"
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

func defaultErrorHandler(status int, err error, c Context) error {
	env := c.Value("env")
	c.Logger().Error(err)
	if env != nil && env.(string) == "production" {
		c.Response().WriteHeader(status)
		c.Response().Write([]byte(prodErrorTmpl))
		return nil
	}
	c.Response().WriteHeader(status)

	msg := fmt.Sprintf("%+v", err)
	ct := c.Request().Header.Get("Content-Type")
	switch strings.ToLower(ct) {
	case "application/json", "text/json", "json":
		err = json.NewEncoder(c.Response()).Encode(map[string]interface{}{
			"error": msg,
			"code":  status,
		})
	case "application/xml", "text/xml", "xml":
	default:
		err := c.Request().ParseForm()
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
	return err
}

var devErrorTmpl = `
<html>
<head>
  <title><%= status %> - ERROR!</title>
  <link rel="stylesheet" href="/assets/application.css" type="text/css" media="all">
  <style>
    .container {
      min-width: 320px;
    }

    body {
      font-family: helvetica;
    }

    table {
      font-size: 14px;
    }

    table.table tbody tr td {
      border-top: 0;
      padding: 10px;
    }

    pre {
      white-space: pre-line;
      margin-bottom: 10px;
      max-height: 275px;
      overflow: scroll;
    }

    header {
      background-color: #ed605e;
      padding: 10px 20px;
      box-sizing: border-box;
    }

    .logo img {
      width: 80px;
    }

    .titles h1 {
      font-size: 30px;
      font-weight: 300;
      color: white;
      margin: 24px 0;
    }

    .content h3 {
      color: gray;
      margin: 25px 0;
    }

    .centered {
      text-align: center;
    }

    .foot {
      padding: 5px 0 20px;
      text-align: right;
      text-align: right;
      color: #c5c5c5;
      font-weight: 300;
    }

    .foot a {
      color: #8b8b8b;
      text-decoration: underline;
    }

    .centered {
      text-align: center;
    }

    @media all and (max-width: 500px) {
      .titles h1 {
        font-size: 25px;
        margin: 26px 0;
      }
    }

    @media all and (max-width: 530px) {
      .titles h1 {
        font-size: 20px;
        margin: 24px 0;
      }
      .logo {
        padding: 0
      }
      .logo img {
        width: 100%;
        max-width: 80px;
      }
    }
  </style>
</head>

<body>
  <header>
    <div class="container">
      <div class="row">
        <div class="col-md-1 col-sm-2 col-xs-3 logo">
          <a href="/"><img src="https://gobuffalo.io/assets/images/logo_med.png" alt=""></a>
        </div>
        <div class="col-md-10 col-sm-6 col-xs-7 titles">
          <h1>
            <%= status %> - ERROR!
          </h1>
        </div>
      </div>
    </div>
  </header>

  <div class="container content">
    <div class="row">
      <div class="col-md-12">
        <h3>Error Trace</h3>
        <pre><%= error %></pre>

        <h3>Context</h3>
        <pre><%= inspect(context) %></pre>

        <h3>Parameters</h3>
        <pre><%= inspect(params) %></pre>

        <h3>Form</h3>
        <pre><%= inspect(posted_form) %></pre>

        <h3>Routes</h3>
        <table class="table table-striped">
          <thead>
            <tr text-align="left">
              <th class="centered">METHOD</th>
              <th>PATH</th>
              <th>NAME</th>
              <th>HANDLER</th>
            </tr>
          </thead>
          <tbody>

            <%= for (r) in routes { %>
              <tr>
                <td class="centered">
                  <%= r.Method %>
                </td>
                <td>
                  <%= r.Path %>
                </td>
                <td>
                  <%= r.PathName %>
                </td>
                <td><code><%= r.HandlerName %></code></td>
              </tr>
            <% } %>

          </tbody>
        </table>
      </div>
    </div>
    <div class="foot"> <span> Powered by <a href="http://gobuffalo.io/">gobuffalo.io</a></span></div>
  </div>
</body>
</html>
`
var prodErrorTmpl = `
<h1>We're Sorry!</h1>
<p>
It looks like something went wrong! Don't worry, we are aware of the problem and are looking into it.
</p>
<p>
Sorry if this has caused you any problems. Please check back again later.
</p>
`
