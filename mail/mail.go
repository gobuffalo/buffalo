// Package mail provides email sending functionality for Buffalo applications.
// It supports SMTP delivery with customizable configuration including TLS/SSL,
// authentication, and batch sending capabilities.
//
// Portions of the SMTP implementation are derived from the go-mail/mail project
// (https://github.com/go-mail/mail) under the MIT License.
//
// TODO: Properly encode filenames for non-ASCII characters.
// TODO: Properly encode email addresses for non-ASCII characters.
// TODO: Test embedded files and attachments for their existence before sending.
// TODO: Allow supplying an io.Reader when embedding and attaching files.
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
