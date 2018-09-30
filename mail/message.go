package mail

import (
	"context"
	"io"
	"sync"

	"bytes"

	"github.com/gobuffalo/buffalo/render"
)

//Message represents an Email message
type Message struct {
	Context     context.Context
	From        string
	To          []string
	CC          []string
	Bcc         []string
	Subject     string
	Headers     map[string]string
	Data        render.Data
	Bodies      []Body
	Attachments []Attachment
	moot        *sync.RWMutex
}

func (m *Message) merge(data render.Data) render.Data {
	d := render.Data{}
	m.moot.Lock()
	for k, v := range m.Data {
		d[k] = v
	}
	m.moot.Unlock()
	for k, v := range data {
		d[k] = v
	}
	return d
}

// AddBody the message by receiving a renderer and rendering data, first message will be
// used as the main message Body rest of them will be passed as alternative bodies on the
// email message
func (m *Message) AddBody(r render.Renderer, data render.Data) error {
	buf := bytes.NewBuffer([]byte{})
	err := r.Render(buf, m.merge(data))

	if err != nil {
		return err
	}

	m.Bodies = append(m.Bodies, Body{
		Content:     buf.String(),
		ContentType: r.ContentType(),
	})

	return nil
}

// AddBodies Allows to add multiple bodies to the message, it returns errors that
// could happen in the rendering.
func (m *Message) AddBodies(data render.Data, renderers ...render.Renderer) error {
	for _, r := range renderers {
		err := m.AddBody(r, data)
		if err != nil {
			return err
		}
	}

	return nil
}

//AddAttachment adds the attachment to the list of attachments the Message has.
func (m *Message) AddAttachment(name, contentType string, r io.Reader) error {
	m.Attachments = append(m.Attachments, Attachment{
		Name:        name,
		ContentType: contentType,
		Reader:      r,
		Embedded:    false,
	})

	return nil
}

//AddEmbedded adds the attachment to the list of attachments
// the Message has and uses inline instead of attachement property.
func (m *Message) AddEmbedded(name string, r io.Reader) error {
	m.Attachments = append(m.Attachments, Attachment{
		Name:     name,
		Reader:   r,
		Embedded: true,
	})

	return nil
}

// SetHeader sets the heder field and value for the message
func (m *Message) SetHeader(field, value string) {
	m.Headers[field] = value
}
