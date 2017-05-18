package buffalo

import (
	"reflect"
	"runtime"
	"strings"
)

// MiddlewareFunc defines the interface for a piece of Buffalo
// Middleware.
/*
	func DoSomething(next Handler) Handler {
		return func(c Context) error {
			// do something before calling the next handler
			err := next(c)
			// do something after call the handler
			return err
		}
	}
*/
type MiddlewareFunc func(Handler) Handler

// Use the specified Middleware for the App.
// When defined on an `*App` the specified middleware will be
// inherited by any `Group` calls that are made on that on
// the App.
func (a *App) Use(mw ...MiddlewareFunc) {
	a.Middleware.Use(mw...)
}

// MiddlewareStack manages the middleware stack for an App/Group.
type MiddlewareStack struct {
	stack []MiddlewareFunc
	skips map[string]bool
}

func (ms *MiddlewareStack) clone() *MiddlewareStack {
	n := newMiddlewareStack()
	for _, s := range ms.stack {
		n.stack = append(n.stack, s)
	}
	for k, v := range ms.skips {
		n.skips[k] = v
	}
	return n
}

// Clear wipes out the current middleware stack for the App/Group,
// any middleware previously defined will be removed leaving an empty
// middleware stack.
func (ms *MiddlewareStack) Clear() {
	ms.stack = []MiddlewareFunc{}
	ms.skips = map[string]bool{}
}

// Use the specified Middleware for the App.
// When defined on an `*App` the specified middleware will be
// inherited by any `Group` calls that are made on that on
// the App.
func (ms *MiddlewareStack) Use(mw ...MiddlewareFunc) {
	ms.stack = append(ms.stack, mw...)
}

// Skip a specified piece of middleware the specified Handlers.
// This is useful for things like wrapping your application in an
// authorization middleware, but skipping it for things the home
// page, the login page, etc...
/*
	a.Middleware.Skip(Authorization, HomeHandler, LoginHandler, RegistrationHandler)
*/
// NOTE: When skipping Resource handlers, you need to first declare your
// resource handler as a type of buffalo.Resource for the Skip function to
// properly recognize and match it.
/*
	// Works:
	var cr Resource
	cr = &carsResource{&buffaloBaseResource{}}
	g = a.Resource("/cars", cr)
	g.Use(SomeMiddleware)
	g.Middleware.Skip(SomeMiddleware, cr.Show)

	// Doesn't Work:
	cr := &carsResource{&buffaloBaseResource{}}
	g = a.Resource("/cars", cr)
	g.Use(SomeMiddleware)
	g.Middleware.Skip(SomeMiddleware, cr.Show)
*/
func (ms *MiddlewareStack) Skip(mw MiddlewareFunc, handlers ...Handler) {
	for _, h := range handlers {
		key := funcKey(mw, h)
		ms.skips[key] = true
	}
}

// Replace a piece of middleware with another piece of middleware. Great for
// testing.
func (ms *MiddlewareStack) Replace(mw1 MiddlewareFunc, mw2 MiddlewareFunc) {
	m1k := funcKey(mw1)
	stack := []MiddlewareFunc{}
	for _, mw := range ms.stack {
		if funcKey(mw) == m1k {
			stack = append(stack, mw2)
		} else {
			stack = append(stack, mw)
		}
	}
	ms.stack = stack
}

func (ms *MiddlewareStack) handler(h Handler) Handler {
	if len(ms.stack) > 0 {
		mh := func(_ Handler) Handler {
			return h
		}

		tstack := []MiddlewareFunc{mh}

		sl := len(ms.stack) - 1
		for i := sl; i >= 0; i-- {
			mw := ms.stack[i]
			key := funcKey(mw, h)
			if !ms.skips[key] {
				tstack = append(tstack, mw)
			}
		}

		for _, mw := range tstack {
			h = mw(h)
		}
		return h
	}
	return h
}

func newMiddlewareStack(mws ...MiddlewareFunc) *MiddlewareStack {
	return &MiddlewareStack{
		stack: mws,
		skips: map[string]bool{},
	}
}

func funcKey(funcs ...interface{}) string {
	names := []string{}
	for _, f := range funcs {
		rv := reflect.ValueOf(f)
		ptr := rv.Pointer()
		if n, ok := keyMap[ptr]; ok {
			names = append(names, n)
			continue
		}
		fnc := runtime.FuncForPC(ptr)
		n := fnc.Name()

		n = strings.Replace(n, "-fm", "", 1)
		n = strings.Replace(n, "(", "", 1)
		n = strings.Replace(n, ")", "", 1)
		keyMap[ptr] = n
		names = append(names, n)
	}
	return strings.Join(names, "/")
}

func setFuncKey(f interface{}, name string) {
	rv := reflect.ValueOf(f)
	ptr := rv.Pointer()
	keyMap[ptr] = name
}

var keyMap = map[uintptr]string{}
