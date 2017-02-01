package buffalo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gobuffalo/velvet"
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

func defaultErrorHandler(status int, err error, c Context) error {
	env := c.Value("env")
	if env != nil && env.(string) == "production" {
		c.Response().WriteHeader(status)
		c.Response().Write([]byte(prodErrorTmpl))
		return nil
	}
	c.Logger().Error(err)
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
		data := map[string]interface{}{
			"routes": c.Value("routes"),
			"error":  msg,
			"status": status,
			"data":   c.Data(),
		}
		ctx := velvet.NewContextWith(data)
		t, err := velvet.Render(devErrorTmpl, ctx)
		if err != nil {
			return errors.WithStack(err)
		}
		res := c.Response()
		res.WriteHeader(404)
		_, err = res.Write([]byte(t))
		return err
	}
	return err
}

var devErrorTmpl = `
<html>
<head>
	<title>{{status}} - ERROR!</title>
	<style>
		body {
			font-family: helvetica;
		}
		table {
			width: 100%;
		}
		th {
			text-align: left;
		}
		tr:nth-child(even) {
		  background-color: #dddddd;
		}
		td {
			margin: 0px;
			padding: 10px;
		}
		pre {
			display: block;
			padding: 9.5px;
			margin: 0 0 10px;
			font-size: 13px;
			line-height: 1.42857143;
			color: #333;
			word-break: break-all;
			word-wrap: break-word;
			background-color: #f5f5f5;
			border: 1px solid #ccc;
			border-radius: 4px;
		}
	</style>
</head>
<body>
<h1>{{status}} - ERROR!</h1>
<pre>{{error}}</pre>
<hr>
<h3>Context</h3>
<pre>{{#each data as |k v|}}
{{inspect k}}: {{inspect v}}
{{/each}}</pre>
<hr>
<h3>Routes</h3>
<table id="buffalo-routes-table">
	<thead>
		<tr>
			<th>METHOD</th>
			<th>PATH</th>
			<th>HANDLER</th>
		</tr>
	</thead>
	<tbody>
		{{#each routes as |route|}}
			<tr>
				<td>{{route.Method}}</td>
				<td>{{route.Path}}</td>
				<td><code>{{route.HandlerName}}</code></td>
			</tr>
		{{/each}}
	</tbody>
</table>
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
