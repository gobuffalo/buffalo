package mail

import (
	"io"

	"bytes"

	"github.com/gobuffalo/buffalo/render"
)

//Message represents an Email message
type Message struct {
	From    string
	To      []string
	CC      []string
	Bcc     []string
	Subject string
	Headers map[string]string

	Bodies      []Body
	Attachments []Attachment
}

// Body represents one of the bodies in the Message could be main or alternative
type Body struct {
	Content     string
	ContentType string
}

// Attachment are files added into a email message
type Attachment struct {
	Name        string
	Reader      io.Reader
	ContentType string
}

// AddBody the message by receiving a renderer and rendering data, first message will be
// used as the main message Body rest of them will be passed as alternative bodies on the
// email message
func (m *Message) AddBody(r render.Renderer, data render.Data) error {
	buf := bytes.NewBuffer([]byte{})
	err := r.Render(buf, data)

	if err != nil {
		return err
	}

	m.Bodies = append(m.Bodies, Body{
		Content:     string(buf.Bytes()),
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
	})

	return nil
}

// SetHeader sets the heder field and value for the message
func (m *Message) SetHeader(field, value string) {
	m.Headers[field] = value
}

//NewMessage Builds a new message.
func NewMessage() Message {
	return Message{Headers: map[string]string{}}
}
