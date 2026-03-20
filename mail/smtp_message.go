// Portions of this code are derived from the go-mail/mail project.
// https://github.com/go-mail/mail (MIT License)

package mail

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"time"
)

// smtpMessage represents an email for SMTP transmission.
type smtpMessage struct {
	header      smtpHeader
	parts       []*part
	attachments []*file
	embedded    []*file
	charset     string
	encoding    encoding
	hEncoder    mimeEncoder
	buf         bytes.Buffer
	boundary    string
}

type smtpHeader map[string][]string

type part struct {
	contentType string
	copier      func(io.Writer) error
	encoding    encoding
}

func newSMTPMessage(settings ...messageSetting) *smtpMessage {
	m := &smtpMessage{
		header:   make(smtpHeader),
		charset:  "UTF-8",
		encoding: encodingQuotedPrintable,
	}

	m.applySettings(settings)

	if m.encoding == encodingBase64 {
		m.hEncoder = bEncoding
	} else {
		m.hEncoder = qEncoding
	}

	return m
}

func (m *smtpMessage) applySettings(settings []messageSetting) {
	for _, s := range settings {
		s(m)
	}
}

type messageSetting func(m *smtpMessage)

type encoding string

const (
	encodingQuotedPrintable encoding = "quoted-printable"
	encodingBase64          encoding = "base64"
	encodingUnencoded       encoding = "8bit"
)

func (m *smtpMessage) setHeader(field string, value ...string) {
	m.encodeHeader(value)
	m.header[field] = value
}

func (m *smtpMessage) encodeHeader(values []string) {
	for i := range values {
		values[i] = m.encodeString(values[i])
	}
}

func (m *smtpMessage) encodeString(value string) string {
	return m.hEncoder.Encode(m.charset, value)
}

func (m *smtpMessage) setBody(contentType, body string, settings ...partSetting) {
	m.setBodyWriter(contentType, newCopier(body), settings...)
}

func (m *smtpMessage) setBodyWriter(contentType string, f func(io.Writer) error, settings ...partSetting) {
	m.parts = []*part{m.newPart(contentType, f, settings)}
}

func (m *smtpMessage) addAlternative(contentType, body string, settings ...partSetting) {
	m.addAlternativeWriter(contentType, newCopier(body), settings...)
}

func newCopier(s string) func(io.Writer) error {
	return func(w io.Writer) error {
		_, err := io.WriteString(w, s)
		return err
	}
}

func (m *smtpMessage) addAlternativeWriter(contentType string, f func(io.Writer) error, settings ...partSetting) {
	m.parts = append(m.parts, m.newPart(contentType, f, settings))
}

func (m *smtpMessage) newPart(contentType string, f func(io.Writer) error, settings []partSetting) *part {
	p := &part{
		contentType: contentType,
		copier:      f,
		encoding:    m.encoding,
	}

	for _, s := range settings {
		s(p)
	}

	return p
}

type partSetting func(*part)

func setPartEncoding(e encoding) partSetting {
	return partSetting(func(p *part) {
		p.encoding = e
	})
}

type file struct {
	Name     string
	Header   map[string][]string
	CopyFunc func(w io.Writer) error
}

func (f *file) setHeader(field, value string) {
	f.Header[field] = []string{value}
}

type fileSetting func(*file)

func setCopyFunc(f func(io.Writer) error) fileSetting {
	return func(fi *file) {
		fi.CopyFunc = f
	}
}

func (m *smtpMessage) attach(filename string, settings ...fileSetting) {
	m.attachments = m.appendFile(m.attachments, fileFromFilename(filename), settings)
}

func (m *smtpMessage) embed(filename string, settings ...fileSetting) {
	m.embedded = m.appendFile(m.embedded, fileFromFilename(filename), settings)
}

func fileFromFilename(name string) *file {
	return &file{
		Name:   filepath.Base(name),
		Header: make(map[string][]string),
		CopyFunc: func(w io.Writer) error {
			h, err := os.Open(name)
			if err != nil {
				return err
			}
			if _, err := io.Copy(w, h); err != nil {
				h.Close()
				return err
			}
			return h.Close()
		},
	}
}

func (m *smtpMessage) formatDate(date time.Time) string {
	return date.Format(time.RFC1123Z)
}

func (m *smtpMessage) appendFile(list []*file, f *file, settings []fileSetting) []*file {
	for _, s := range settings {
		s(f)
	}

	if list == nil {
		return []*file{f}
	}

	return append(list, f)
}
