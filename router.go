package buffalo

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"sort"

	"github.com/markbates/inflect"
)

// GET maps an HTTP "GET" request to the path and the specified handler.
func (a *App) GET(p string, h Handler) {
	a.addRoute("GET", p, h)
}

// POST maps an HTTP "POST" request to the path and the specified handler.
func (a *App) POST(p string, h Handler) {
	a.addRoute("POST", p, h)
}

// PUT maps an HTTP "PUT" request to the path and the specified handler.
func (a *App) PUT(p string, h Handler) {
	a.addRoute("PUT", p, h)
}

// DELETE maps an HTTP "DELETE" request to the path and the specified handler.
func (a *App) DELETE(p string, h Handler) {
	a.addRoute("DELETE", p, h)
}

// HEAD maps an HTTP "HEAD" request to the path and the specified handler.
func (a *App) HEAD(p string, h Handler) {
	a.addRoute("HEAD", p, h)
}

// OPTIONS maps an HTTP "OPTIONS" request to the path and the specified handler.
func (a *App) OPTIONS(p string, h Handler) {
	a.addRoute("OPTIONS", p, h)
}

// PATCH maps an HTTP "PATCH" request to the path and the specified handler.
func (a *App) PATCH(p string, h Handler) {
	a.addRoute("PATCH", p, h)
}

// ServeFiles maps an path to a directory on disk to serve static files.
// Useful for JavaScript, images, CSS, etc...
/*
	a.ServeFiles("/assets", http.Dir("path/to/assets"))
*/
func (a *App) ServeFiles(p string, root http.FileSystem) {
	a.router.PathPrefix(p).Handler(http.StripPrefix(p, http.FileServer(root)))
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
	base := filepath.Base(p)
	single := inflect.Singularize(base)
	g := a.Group(p)
	p = "/"
	spath := filepath.Join(p, fmt.Sprintf("{%s_id}", single))
	g.GET(p, r.List)
	g.GET(filepath.Join(p, "new"), r.New)
	g.GET(filepath.Join(spath), r.Show)
	g.GET(filepath.Join(spath, "edit"), r.Edit)
	g.POST(p, r.Create)
	g.PUT(filepath.Join(spath), r.Update)
	g.DELETE(filepath.Join(spath), r.Destroy)
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
func (a *App) Group(path string) *App {
	g := New(a.Options)
	g.prefix = path
	g.router = a.router
	g.Middleware = newMiddlewareStack()
	g.Middleware.skips = a.Middleware.skips
	g.Middleware.stack = a.Middleware.stack
	g.root = a
	if a.root != nil {
		g.root = a.root
	}
	return g
}

func (a *App) addRoute(method string, url string, h Handler) {
	a.moot.Lock()
	defer a.moot.Unlock()

	url = path.Join(a.prefix, url)
	hs := funcKey(h)
	routes := a.Routes()
	routes = append(routes, route{
		Method:      method,
		Path:        url,
		HandlerName: hs,
	})

	sort.Sort(routes)
	if a.root != nil {
		a.root.routes = routes
	} else {
		a.routes = routes
	}

	a.router.Handle(url, a.handlerToHandler(h)).Methods(method)
}
