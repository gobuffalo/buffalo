package vm

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/mattn/anko/parser"
)

// Env provides interface to run VM. This mean function scope and blocked-scope.
// If stack goes to blocked-scope, it will make new Env.
type Env struct {
	name      string
	env       map[string]reflect.Value
	typ       map[string]reflect.Type
	parent    *Env
	interrupt *bool
	sync.RWMutex
}

// NewEnv creates new global scope.
func NewEnv() *Env {
	b := false

	return &Env{
		env:       make(map[string]reflect.Value),
		typ:       make(map[string]reflect.Type),
		parent:    nil,
		interrupt: &b,
	}
}

// NewEnv creates new child scope.
func (e *Env) NewEnv() *Env {
	return &Env{
		env:       make(map[string]reflect.Value),
		typ:       make(map[string]reflect.Type),
		parent:    e,
		name:      e.name,
		interrupt: e.interrupt,
	}
}

func NewPackage(n string) *Env {
	b := false

	return &Env{
		env:       make(map[string]reflect.Value),
		typ:       make(map[string]reflect.Type),
		parent:    nil,
		name:      n,
		interrupt: &b,
	}
}

func (e *Env) NewPackage(n string) *Env {
	return &Env{
		env:       make(map[string]reflect.Value),
		typ:       make(map[string]reflect.Type),
		parent:    e,
		name:      n,
		interrupt: e.interrupt,
	}
}

// Destroy deletes current scope.
func (e *Env) Destroy() {
	e.Lock()
	defer e.Unlock()

	if e.parent == nil {
		return
	}
	for k, v := range e.parent.env {
		if v.IsValid() && v.Interface() == e {
			delete(e.parent.env, k)
		}
	}
	e.parent = nil
	e.env = nil
}

// NewModule creates new module scope as global.
func (e *Env) NewModule(n string) *Env {
	m := &Env{
		env:    make(map[string]reflect.Value),
		parent: e,
		name:   n,
	}
	e.Define(n, m)
	return m
}

// SetName sets a name of the scope. This means that the scope is module.
func (e *Env) SetName(n string) {
	e.Lock()
	e.name = n
	e.Unlock()
}

// GetName returns module name.
func (e *Env) GetName() string {
	e.RLock()
	defer e.RUnlock()

	return e.name
}

// Addr returns pointer value which specified symbol. It goes to upper scope until
// found or returns error.
func (e *Env) Addr(k string) (reflect.Value, error) {
	e.RLock()
	defer e.RUnlock()

	if v, ok := e.env[k]; ok {
		return v.Addr(), nil
	}
	if e.parent == nil {
		return NilValue, fmt.Errorf("Undefined symbol '%s'", k)
	}
	return e.parent.Addr(k)
}

// Type returns type which specified symbol. It goes to upper scope until
// found or returns error.
func (e *Env) Type(k string) (reflect.Type, error) {
	e.RLock()
	defer e.RUnlock()

	if v, ok := e.typ[k]; ok {
		return v, nil
	}
	if e.parent == nil {
		return NilType, fmt.Errorf("Undefined type '%s'", k)
	}
	return e.parent.Type(k)
}

// Get returns value which specified symbol. It goes to upper scope until
// found or returns error.
func (e *Env) Get(k string) (reflect.Value, error) {
	e.RLock()
	defer e.RUnlock()

	if v, ok := e.env[k]; ok {
		return v, nil
	}
	if e.parent == nil {
		return NilValue, fmt.Errorf("Undefined symbol '%s'", k)
	}
	return e.parent.Get(k)
}

// Set modifies value which specified as symbol. It goes to upper scope until
// found or returns error.
func (e *Env) Set(k string, v interface{}) error {
	e.Lock()
	defer e.Unlock()

	if _, ok := e.env[k]; ok {
		val, ok := v.(reflect.Value)
		if !ok {
			val = reflect.ValueOf(v)
		}
		e.env[k] = val
		return nil
	}
	if e.parent == nil {
		return fmt.Errorf("Unknown symbol '%s'", k)
	}
	return e.parent.Set(k, v)
}

// DefineGlobal defines symbol in global scope.
func (e *Env) DefineGlobal(k string, v interface{}) error {
	if e.parent == nil {
		return e.Define(k, v)
	}
	return e.parent.DefineGlobal(k, v)
}

// DefineType defines type which specifis symbol in global scope.
func (e *Env) DefineType(k string, t interface{}) error {
	if strings.Contains(k, ".") {
		return fmt.Errorf("Unknown symbol '%s'", k)
	}
	global := e
	keys := []string{k}

	e.RLock()
	for global.parent != nil {
		if global.name != "" {
			keys = append(keys, global.name)
		}
		global = global.parent
	}
	e.RUnlock()

	for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
		keys[i], keys[j] = keys[j], keys[i]
	}

	typ, ok := t.(reflect.Type)
	if !ok {
		typ = reflect.TypeOf(t)
	}

	global.Lock()
	global.typ[strings.Join(keys, ".")] = typ
	global.Unlock()

	return nil
}

// Define defines symbol in current scope.
func (e *Env) Define(k string, v interface{}) error {
	if strings.Contains(k, ".") {
		return fmt.Errorf("Unknown symbol '%s'", k)
	}
	val, ok := v.(reflect.Value)
	if !ok {
		val = reflect.ValueOf(v)
	}

	e.Lock()
	e.env[k] = val
	e.Unlock()

	return nil
}

// String return the name of current scope.
func (e *Env) String() string {
	e.RLock()
	defer e.RUnlock()

	return e.name
}

// Dump show symbol values in the scope.
func (e *Env) Dump() {
	e.RLock()
	for k, v := range e.env {
		fmt.Printf("%v = %#v\n", k, v)
	}
	e.RUnlock()
}

// Execute parses and runs source in current scope.
func (e *Env) Execute(src string) (reflect.Value, error) {
	stmts, err := parser.ParseSrc(src)
	if err != nil {
		return NilValue, err
	}
	return Run(stmts, e)
}
