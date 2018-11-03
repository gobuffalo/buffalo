package mail

import (
	"context"
	"sync"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

// NewMessage builds a new message.
func NewMessage() Message {
	return Message{
		Context: context.Background(),
		Headers: map[string]string{},
		Data:    render.Data{},
		moot:    &sync.RWMutex{},
	}
}

// NewFromData builds a new message with raw template data given
func NewFromData(data render.Data) Message {
	d := render.Data{}
	for k, v := range data {
		d[k] = v
	}
	m := NewMessage()
	m.Data = d
	return m
}

// New builds a new message with the current buffalo.Context
func New(c buffalo.Context) Message {
	m := NewFromData(c.Data())
	m.Context = c
	return m
}
