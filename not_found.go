package buffalo

import (
	"encoding/json"
	"net/http"

	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
)

// NotFoundHandler is the default ErrorHandler for 404 responses.
// In development mode it attempts to return useful debugging
// information. In production it defaults to use the http.NotFound
// handler.
func NotFoundHandler(status int, err error, c Context) error {
	env := c.Value("env")
	req := c.Request()
	res := c.Response()
	if env != nil && env.(string) == "production" {
		http.NotFound(res, req)
		return nil
	}
	data := map[string]interface{}{
		"routes": c.Value("routes"),
		"method": req.Method,
		"path":   req.URL.String(),
		"error":  err.Error(),
	}
	ct := req.Header.Get("Content-Type")
	if ct == "application/json" {
		res.WriteHeader(404)
		return json.NewEncoder(res).Encode(data)
	}
	ctx := plush.NewContextWith(data)
	t, err := plush.Render(htmlNotFound, ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	res.WriteHeader(404)
	_, err = res.Write([]byte(t))
	return err
}

var htmlNotFound = `
<html>
<head>
	<title>404 PAGE NOT FOUND</title>
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
	</style>
</head>
<body>
<h1>404 Page Not Found!</h1>
<h3>Could not find path <code>[<%= method %>] <%= path %></code></h3>
<hr>
<table id="buffalo-routes-table">
	<thead>
		<tr>
			<th>METHOD</th>
			<th>PATH</th>
			<th>NAME</th>
			<th>HANDLER</th>
		</tr>
	</thead>
	<tbody>
		<%= for (route) in routes { %>
			<tr>
				<td><%= route.Method %></td>
				<td><%= route.Path %></td>
				<td><%= route.PathName %></td>
				<td><code><%= route.HandlerName %></code></td>
			</tr>
		<% } %>
	</tbody>
</table>
<%= if (error) { %>
<hr>
<h2>Error</h2>
<pre><%= error %></pre>
<% } %>
</body>
</html>
`
