package buffalo

func (a *App) Routes() RouteList {
	if a.root != nil {
		return a.root.routes
	}
	return a.routes
}

type route struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	HandlerName string `json:"handler"`
}

type RouteList []route

func (a RouteList) Len() int      { return len(a) }
func (a RouteList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a RouteList) Less(i, j int) bool {
	x := a[i].Method + a[i].Path
	y := a[j].Method + a[j].Path
	return x < y
}
