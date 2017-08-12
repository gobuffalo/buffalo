package plush

import (
	"context"
	"sync"
)

var _ context.Context = &Context{}

// Context holds all of the data for the template that is being rendered.
type Context struct {
	context.Context
	data  map[string]interface{}
	outer *Context
	moot  *sync.Mutex
}

// New context containing the current context. Values set on the new context
// will not be set onto the original context, however, the original context's
// values will be available to the new context.
func (c *Context) New() *Context {
	cc := NewContext()
	cc.outer = c
	return cc
}

// Set a value onto the context
func (c *Context) Set(key string, value interface{}) {
	c.moot.Lock()
	defer c.moot.Unlock()
	c.data[key] = value
}

// Value from the context, or it's parent's context if one exists.
func (c *Context) Value(key interface{}) interface{} {
	if s, ok := key.(string); ok {
		if v, ok := c.data[s]; ok {
			return v
		}
		if c.outer != nil {
			return c.outer.Value(s)
		}
	}
	return c.Context.Value(key)
}

// Has checks the existence of the key in the context.
func (c *Context) Has(key string) bool {
	return c.Value(key) != nil
}

func (c *Context) export() map[string]interface{} {
	m := map[string]interface{}{}
	if c.outer != nil {
		for k, v := range c.outer.export() {
			m[k] = v
		}
	}
	for k, v := range c.data {
		m[k] = v
	}

	return m
}

// NewContext returns a fully formed context ready to go
func NewContext() *Context {
	return &Context{
		Context: context.Background(),
		data:    map[string]interface{}{},
		outer:   nil,
		moot:    &sync.Mutex{},
	}
}

// NewContextWith returns a fully formed context using the data
// provided.
func NewContextWith(data map[string]interface{}) *Context {
	c := NewContext()
	c.data = data
	return c
}

// NewContextWithContext returns a new plush.Context given another context
func NewContextWithContext(ctx context.Context) *Context {
	c := NewContext()
	c.Context = ctx
	return c
}
