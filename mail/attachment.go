package mail

import "io"

// Attachment are files added into a email message
type Attachment struct {
	Name        string
	Reader      io.Reader
	ContentType string
	Embedded    bool
}
