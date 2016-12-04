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
// authorization middleare, but skipping it for things the home
// page, the login page, etc...
/*
	a.Middleware.Skip(Authorization, HomeHandler, LoginHandler, RegistrationHandler)
*/
func (ms *MiddlewareStack) Skip(mw MiddlewareFunc, handlers ...Handler) {
	for _, h := range handlers {
		ms.skips[funcKey(mw, h)] = true
	}
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
			if !ms.skips[funcKey(mw, h)] {
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
		n := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		names = append(names, n)
	}
	return strings.Join(names, "/")
}
