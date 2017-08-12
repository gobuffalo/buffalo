// Package indentwriter implements an io.Writer wrapper that indents
// every non-empty line with specified number of tabs.
package indentwriter

import (
	"bytes"
	"io"
)

type indentWriter struct {
	w      io.Writer
	prefix []byte

	wroteIndent bool
}

// New creates a new indent writer that indents non-empty lines with indent number of tabs.
func New(w io.Writer, indent int) io.Writer {
	return &indentWriter{
		w:      w,
		prefix: bytes.Repeat([]byte{'\t'}, indent),
	}
}

func (iw *indentWriter) Write(p []byte) (n int, err error) {
	for i, b := range p {
		if b == '\n' {
			iw.wroteIndent = false
		} else {
			if !iw.wroteIndent {
				_, err = iw.w.Write(iw.prefix)
				if err != nil {
					return n, err
				}
				iw.wroteIndent = true
			}
		}
		_, err = iw.w.Write(p[i : i+1])
		if err != nil {
			return n, err
		}
		n++
	}
	return len(p), nil
}
