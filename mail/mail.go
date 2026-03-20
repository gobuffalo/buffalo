package mail

import (
	"context"
	"maps"
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
	m := NewMessage()
	m.Data = maps.Clone(data)
	return m
}

// New builds a new message with the current buffalo.Context
func New(c buffalo.Context) Message {
	m := NewFromData(c.Data())
	m.Context = c
	return m
}
