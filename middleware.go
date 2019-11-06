package buffalo

import (
	"reflect"
	"runtime"
	"strings"
	"sync"
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

const funcKeyDelimeter = ":"

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

func (ms MiddlewareStack) String() string {
	s := []string{}
	for _, m := range ms.stack {
		s = append(s, funcKey(m))
	}

	return strings.Join(s, "\n")
}

func (ms *MiddlewareStack) clone() *MiddlewareStack {
	n := newMiddlewareStack()
	n.stack = append(n.stack, ms.stack...)
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

// Remove the specified Middleware(s) for the App/group. This is useful when 
// the middleware will be skipped by the entire group.
/*
	a.Middleware.Remove(Authorization)
*/
*/
func (ms *MiddlewareStack) Remove(mws ...MiddlewareFunc) {
	result := []MiddlewareFunc{}

base:
	for _, existing := range ms.stack {
		for _, banned := range mws {
			if funcKey(existing) == funcKey(banned) {
				continue base
			}
		}

		result = append(result, existing)
	}

	ms.stack = result

}

// Skip a specified piece of middleware the specified Handlers.
// This is useful for things like wrapping your application in an
// authorization middleware, but skipping it for things the home
// page, the login page, etc...
/*
	a.Middleware.Skip(Authorization, HomeHandler, LoginHandler, RegistrationHandler)
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

func (ms *MiddlewareStack) handler(info RouteInfo) Handler {
	h := info.Handler
	if len(ms.stack) > 0 {
		mh := func(_ Handler) Handler {
			return h
		}

		tstack := []MiddlewareFunc{mh}

		sl := len(ms.stack) - 1
		for i := sl; i >= 0; i-- {
			mw := ms.stack[i]
			key := funcKey(mw, info)
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
		if n, ok := f.(RouteInfo); ok {
			names = append(names, n.HandlerName)
			continue
		}
		rv := reflect.ValueOf(f)
		ptr := rv.Pointer()
		keyMapMutex.Lock()
		if n, ok := keyMap[ptr]; ok {
			keyMapMutex.Unlock()
			names = append(names, n)
			continue
		}
		keyMapMutex.Unlock()
		n := ptrName(ptr)
		keyMapMutex.Lock()
		keyMap[ptr] = n
		keyMapMutex.Unlock()
		names = append(names, n)
	}
	return strings.Join(names, funcKeyDelimeter)
}

func ptrName(ptr uintptr) string {
	fnc := runtime.FuncForPC(ptr)
	n := fnc.Name()

	n = strings.Replace(n, "-fm", "", 1)
	n = strings.Replace(n, "(", "", 1)
	n = strings.Replace(n, ")", "", 1)
	return n
}

func setFuncKey(f interface{}, name string) {
	rv := reflect.ValueOf(f)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	ptr := rv.Pointer()
	keyMapMutex.Lock()
	keyMap[ptr] = name
	keyMapMutex.Unlock()
}

var keyMap = map[uintptr]string{}
var keyMapMutex = sync.Mutex{}
