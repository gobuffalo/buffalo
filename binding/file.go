package binding

import (
	"mime/multipart"
	"net/http"
	"reflect"

	"github.com/pkg/errors"
)

// MaxFileSize can be used to set the maximum file size for multipart
// requests. See https://golang.org/pkg/net/http/#Request.ParseMultipartForm
// for more information on how this impacts file uploads.
var MaxFileSize int64 = 5 * 1024 * 1024

func init() {
	sb := func(req *http.Request, i interface{}) error {
		err := req.ParseMultipartForm(MaxFileSize)
		if err != nil {
			return errors.WithStack(err)
		}
		if err := decoder.Decode(req.Form, i); err != nil {
			return errors.WithStack(err)
		}

		form := req.MultipartForm.File
		if len(form) == 0 {
			return nil
		}

		ri := reflect.Indirect(reflect.ValueOf(i))
		for n, _ := range form {
			f := ri.FieldByName(n)
			if _, ok := f.Interface().(File); !ok {
				continue
			}
			mf, mh, err := req.FormFile(n)
			if err != nil {
				return errors.WithStack(err)
			}
			f.Set(reflect.ValueOf(File{
				File:       mf,
				FileHeader: mh,
			}))
		}

		return nil
	}
	binders["multipart/form-data"] = sb
}

// File holds information regarding an uploaded file
type File struct {
	multipart.File
	*multipart.FileHeader
}

// Valid if there is an actual uploaded file
func (f File) Valid() bool {
	if f.File == nil {
		return false
	}
	return true
}

func (f File) String() string {
	if f.File == nil {
		return ""
	}
	return f.Filename
}
