package buffalo

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/pkg/errors"
)

func (a *App) notFound() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if a.Env == "development" {
			err := func() error {
				routes := a.Routes()
				data := map[string]interface{}{
					"routes": routes,
					"method": req.Method,
					"path":   req.URL.String(),
				}
				switch req.Header.Get("Content-Type") {
				case "application/json":
					res.WriteHeader(404)
					return json.NewEncoder(res).Encode(data)
				default:
					t, err := template.New("not-found").Parse(htmlNotFound)
					if err != nil {
						res.WriteHeader(500)
						err = errors.WithStack(err)
						res.Write([]byte(err.Error()))
						return err
					}
					res.WriteHeader(404)
					return t.Execute(res, data)
				}
			}()
			if err != nil {
				a.Logger.Error(err)
			}
			return
		}
		http.NotFound(res, req)
	})
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
<h3>Could not find path <code>[{{.method}}] {{.path}}</code></h3>
<hr>
<table id="buffalo-routes-table">
	<thead>
		<tr>
			<th>METHOD</th>
			<th>PATH</th>
			<th>HANDLER</th>
		</tr>
	</thead>
	<tbody>
		{{range .routes}}
			<tr>
				<td>{{.Method}}</td>
				<td>{{.Path}}</td>
				<td><code>{{.HandlerName}}</code></td>
			</tr>
		{{end}}
	</tbody>
</table>
</body>
</html>
`
