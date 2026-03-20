// Portions of this code are derived from the go-mail/mail project.
// https://github.com/go-mail/mail (MIT License)

package mail

import (
	"mime"
	"mime/quotedprintable"
	"strings"
)

var newQPWriter = quotedprintable.NewWriter

type mimeEncoder struct {
	mime.WordEncoder
}

var (
	bEncoding     = mimeEncoder{mime.BEncoding}
	qEncoding     = mimeEncoder{mime.QEncoding}
	lastIndexByte = strings.LastIndexByte
)
