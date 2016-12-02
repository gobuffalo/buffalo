package buffalo

import (
	"net/http"
	"path"
	"sort"
)

func (a *App) GET(p string, h Handler) {
	a.addRoute("GET", p, h)
}

func (a *App) POST(p string, h Handler) {
	a.addRoute("POST", p, h)
}

func (a *App) PUT(p string, h Handler) {
	a.addRoute("PUT", p, h)
}

func (a *App) DELETE(p string, h Handler) {
	a.addRoute("DELETE", p, h)
}

func (a *App) HEAD(p string, h Handler) {
	a.addRoute("HEAD", p, h)
}

func (a *App) OPTIONS(p string, h Handler) {
	a.addRoute("OPTIONS", p, h)
}

func (a *App) PATCH(p string, h Handler) {
	a.addRoute("PATCH", p, h)
}

func (a *App) ServeFiles(p string, root http.FileSystem) {
	a.router.ServeFiles(path.Join(p, "*filepath"), root)
}

func (a *App) ANY(p string, h Handler) {
	a.GET(p, h)
	a.POST(p, h)
	a.PUT(p, h)
	a.PATCH(p, h)
	a.HEAD(p, h)
	a.OPTIONS(p, h)
	a.DELETE(p, h)
}

func (a *App) Group(path string) *App {
	g := New(a.Options)
	g.prefix = path
	g.router = a.router
	g.middlewareStack = a.middlewareStack
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

	a.router.Handle(method, url, a.handlerToHandler(h))
}
