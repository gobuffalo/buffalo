package binding

import (
	"mime/multipart"
)

// File holds information regarding an uploaded file
type File struct {
	multipart.File
	*multipart.FileHeader
}

// Valid if there is an actual uploaded file
func (f File) Valid() bool {
	return f.File != nil
}

func (f File) String() string {
	if f.File == nil {
		return ""
	}
	return f.Filename
}
