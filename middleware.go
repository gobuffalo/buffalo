package buffalo

import (
	"reflect"
	"runtime"
	"strings"
)

type MiddlewareFunc func(Handler) Handler

func (a *App) Use(mw MiddlewareFunc) {
	stack := a.middlewareStack.stack
	a.middlewareStack.stack = append(stack, mw)
}

func (a *App) Skip(mw MiddlewareFunc, handlers ...Handler) {
	ms := &a.middlewareStack
	ms.skip(mw, handlers...)
}

type middlewareStack struct {
	stack []MiddlewareFunc
	skips map[string]bool
}

func (ms *middlewareStack) skip(mw MiddlewareFunc, handlers ...Handler) {
	for _, h := range handlers {
		ms.skips[funcKey(mw, h)] = true
	}
}

func (ms *middlewareStack) handler(h Handler) Handler {
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

func newMiddlewareStack() middlewareStack {
	return middlewareStack{
		stack: []MiddlewareFunc{},
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
