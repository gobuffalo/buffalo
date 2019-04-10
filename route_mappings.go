package buffalo

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/flect/name"
)

const (
	// AssetsAgeVarName is the ENV variable used to specify max age when ServeFiles is used.
	AssetsAgeVarName = "ASSETS_MAX_AGE"
)

// GET maps an HTTP "GET" request to the path and the specified handler.
func (a *App) GET(p string, h Handler) *RouteInfo {
	return a.addRoute("GET", p, h)
}

// POST maps an HTTP "POST" request to the path and the specified handler.
func (a *App) POST(p string, h Handler) *RouteInfo {
	return a.addRoute("POST", p, h)
}

// PUT maps an HTTP "PUT" request to the path and the specified handler.
func (a *App) PUT(p string, h Handler) *RouteInfo {
	return a.addRoute("PUT", p, h)
}

// DELETE maps an HTTP "DELETE" request to the path and the specified handler.
func (a *App) DELETE(p string, h Handler) *RouteInfo {
	return a.addRoute("DELETE", p, h)
}

// HEAD maps an HTTP "HEAD" request to the path and the specified handler.
func (a *App) HEAD(p string, h Handler) *RouteInfo {
	return a.addRoute("HEAD", p, h)
}

// OPTIONS maps an HTTP "OPTIONS" request to the path and the specified handler.
func (a *App) OPTIONS(p string, h Handler) *RouteInfo {
	return a.addRoute("OPTIONS", p, h)
}

// PATCH maps an HTTP "PATCH" request to the path and the specified handler.
func (a *App) PATCH(p string, h Handler) *RouteInfo {
	return a.addRoute("PATCH", p, h)
}

// Redirect from one URL to another URL. Only works for "GET" requests.
func (a *App) Redirect(status int, from, to string) *RouteInfo {
	return a.GET(from, func(c Context) error {
		return c.Redirect(status, to)
	})
}

// Mount mounts a http.Handler (or Buffalo app) and passes through all requests to it.
//
//	func muxer() http.Handler {
//		f := func(res http.ResponseWriter, req *http.Request) {
//			fmt.Fprintf(res, "%s - %s", req.Method, req.URL.String())
//		}
//		mux := mux.NewRouter()
//		mux.HandleFunc("/foo", f).Methods("GET")
//		mux.HandleFunc("/bar", f).Methods("POST")
//		mux.HandleFunc("/baz/baz", f).Methods("DELETE")
//		return mux
//	}
//
//	a.Mount("/admin", muxer())
//
//	$ curl -X DELETE http://localhost:3000/admin/baz/baz
func (a *App) Mount(p string, h http.Handler) {
	prefix := path.Join(a.Prefix, p)
	path := path.Join(p, "{path:.+}")
	a.ANY(path, WrapHandler(http.StripPrefix(prefix, h)))
}

// ServeFiles maps an path to a directory on disk to serve static files.
// Useful for JavaScript, images, CSS, etc...
/*
	a.ServeFiles("/assets", http.Dir("path/to/assets"))
*/
func (a *App) ServeFiles(p string, root http.FileSystem) {
	path := path.Join(a.Prefix, p)
	a.filepaths = append(a.filepaths, path)

	h := stripAsset(path, a.fileServer(root), a)
	a.router.PathPrefix(path).Handler(h)
}

func (a *App) fileServer(fs http.FileSystem) http.Handler {
	fsh := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			eh := a.ErrorHandlers.Get(404)
			eh(404, fmt.Errorf("could not find %s", r.URL.Path), a.newContext(RouteInfo{}, w, r))
			return
		}

		stat, _ := f.Stat()
		maxAge := envy.Get(AssetsAgeVarName, "31536000")
		w.Header().Add("ETag", fmt.Sprintf("%x", stat.ModTime().UnixNano()))
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%s", maxAge))
		fsh.ServeHTTP(w, r)
	})
}

type newable interface {
	New(Context) error
}

type editable interface {
	Edit(Context) error
}

// Resource maps an implementation of the Resource interface
// to the appropriate RESTful mappings. Resource returns the *App
// associated with this group of mappings so you can set middleware, etc...
// on that group, just as if you had used the a.Group functionality.
/*
	a.Resource("/users", &UsersResource{})

	// Is equal to this:

	ur := &UsersResource{}
	g := a.Group("/users")
	g.GET("/", ur.List) // GET /users => ur.List
	g.GET("/new", ur.New) // GET /users/new => ur.New
	g.GET("/{user_id}", ur.Show) // GET /users/{user_id} => ur.Show
	g.GET("/{user_id}/edit", ur.Edit) // GET /users/{user_id}/edit => ur.Edit
	g.POST("/", ur.Create) // POST /users => ur.Create
	g.PUT("/{user_id}", ur.Update) PUT /users/{user_id} => ur.Update
	g.DELETE("/{user_id}", ur.Destroy) DELETE /users/{user_id} => ur.Destroy
*/
func (a *App) Resource(p string, r Resource) *App {

	g := a.Group(p)
	p = "/"

	rv := reflect.ValueOf(r)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	rt := rv.Type()
	rname := fmt.Sprintf("%s.%s", rt.PkgPath(), rt.Name()) + ".%s"

	n := strings.TrimSuffix(rt.Name(), "Resource")
	paramName := name.New(n).ParamID().String()

	type paramKeyable interface {
		ParamKey() string
	}

	if pk, ok := r.(paramKeyable); ok {
		paramName = pk.ParamKey()
	}

	spath := path.Join(p, "{"+paramName+"}")
	setFuncKey(r.List, fmt.Sprintf(rname, "List"))
	g.GET(p, r.List)

	if n, ok := r.(newable); ok {
		setFuncKey(n.New, fmt.Sprintf(rname, "New"))
		g.GET(path.Join(p, "new"), n.New)
	}

	setFuncKey(r.Show, fmt.Sprintf(rname, "Show"))
	g.GET(path.Join(spath), r.Show)

	if n, ok := r.(editable); ok {
		setFuncKey(n.Edit, fmt.Sprintf(rname, "Edit"))
		g.GET(path.Join(spath, "edit"), n.Edit)
	}

	setFuncKey(r.Create, fmt.Sprintf(rname, "Create"))
	g.POST(p, r.Create)
	setFuncKey(r.Update, fmt.Sprintf(rname, "Update"))
	g.PUT(path.Join(spath), r.Update)
	setFuncKey(r.Destroy, fmt.Sprintf(rname, "Destroy"))
	g.DELETE(path.Join(spath), r.Destroy)
	g.Prefix = path.Join(g.Prefix, spath)
	return g
}

// ANY accepts a request across any HTTP method for the specified path
// and routes it to the specified Handler.
func (a *App) ANY(p string, h Handler) {
	a.GET(p, h)
	a.POST(p, h)
	a.PUT(p, h)
	a.PATCH(p, h)
	a.HEAD(p, h)
	a.OPTIONS(p, h)
	a.DELETE(p, h)
}

// Group creates a new `*App` that inherits from it's parent `*App`.
// This is useful for creating groups of end-points that need to share
// common functionality, like middleware.
/*
	g := a.Group("/api/v1")
	g.Use(AuthorizeAPIMiddleware)
	g.GET("/users, APIUsersHandler)
	g.GET("/users/:user_id, APIUserShowHandler)
*/
func (a *App) Group(groupPath string) *App {
	g := New(a.Options)
	g.Prefix = path.Join(a.Prefix, groupPath)
	g.Name = g.Prefix

	g.router = a.router
	g.Middleware = a.Middleware.clone()
	g.ErrorHandlers = a.ErrorHandlers
	g.root = a
	if a.root != nil {
		g.root = a.root
	}
	a.children = append(a.children, g)
	return g
}

// RouteHelpers returns a map of BuildPathHelper() for each route available in the app.
func (a *App) RouteHelpers() map[string]RouteHelperFunc {
	rh := map[string]RouteHelperFunc{}
	for _, route := range a.Routes() {
		cRoute := route
		rh[cRoute.PathName] = cRoute.BuildPathHelper()
	}
	return rh
}

func (a *App) addRoute(method string, url string, h Handler) *RouteInfo {
	a.moot.Lock()
	defer a.moot.Unlock()

	url = path.Join(a.Prefix, url)
	url = a.normalizePath(url)
	name := a.buildRouteName(url)

	hs := funcKey(h)
	r := &RouteInfo{
		Method:      method,
		Path:        url,
		HandlerName: hs,
		Handler:     h,
		App:         a,
		Aliases:     []string{},
	}

	r.MuxRoute = a.router.Handle(url, r).Methods(method)
	r.Name(name)

	routes := a.Routes()
	routes = append(routes, r)
	sort.Sort(routes)

	if a.root != nil {
		a.root.routes = routes
	} else {
		a.routes = routes
	}

	return r
}

//buildRouteName builds a route based on the path passed.
func (a *App) buildRouteName(p string) string {
	if p == "/" || p == "" {
		return "root"
	}

	resultParts := []string{}
	parts := strings.Split(p, "/")

	for index, part := range parts {

		if strings.Contains(part, "{") || part == "" {
			continue
		}

		shouldSingularize := (len(parts) > index+1) && strings.Contains(parts[index+1], "{")
		if shouldSingularize {
			part = flect.Singularize(part)
		}

		if parts[index] == "new" || parts[index] == "edit" {
			resultParts = append([]string{part}, resultParts...)
			continue
		}

		if index > 0 && strings.Contains(parts[index-1], "}") {
			resultParts = append(resultParts, part)
			continue
		}

		resultParts = append(resultParts, part)
	}

	if len(resultParts) == 0 {
		return "unnamed"
	}

	underscore := strings.TrimSpace(strings.Join(resultParts, "_"))
	return name.VarCase(underscore)
}

func stripAsset(path string, h http.Handler, a *App) http.Handler {
	if path == "" {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		up := r.URL.Path
		up = strings.TrimPrefix(up, path)
		up = strings.TrimSuffix(up, "/")

		u, err := url.Parse(up)
		if err != nil {
			eh := a.ErrorHandlers.Get(400)
			eh(400, err, a.newContext(RouteInfo{}, w, r))
			return
		}

		r.URL = u
		h.ServeHTTP(w, r)
	})
}
