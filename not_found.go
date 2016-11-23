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
			routes := a.Routes()
			data := map[string]interface{}{
				"routes": routes,
				"method": req.Method,
				"path":   req.URL.String(),
			}
			switch req.Header.Get("Content-Type") {
			case "application/json":
				res.WriteHeader(404)
				json.NewEncoder(res).Encode(data)
				return
			default:
				t, err := template.New("not-found").Parse(htmlNotFound)
				if err != nil {
					res.WriteHeader(500)
					res.Write([]byte(errors.WithStack(err).Error()))
					return
				}
				res.WriteHeader(404)
				t.Execute(res, data)
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
</head>
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
				<td>{{.HandlerName}}</td>
			</tr>
		{{end}}
	</tbody>
</table>
</html>
`
