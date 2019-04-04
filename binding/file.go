package binding

import (
	"mime/multipart"
	"net/http"
	"reflect"
)

// MaxFileMemory can be used to set the maximum size, in bytes, for files to be
// stored in memory during uploaded for multipart requests.
// See https://golang.org/pkg/net/http/#Request.ParseMultipartForm for more
// information on how this impacts file uploads.
var MaxFileMemory int64 = 5 * 1024 * 1024

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

func init() {
	sb := func(req *http.Request, i interface{}) error {
		err := req.ParseMultipartForm(MaxFileMemory)
		if err != nil {
			return err
		}
		if err := decoder.Decode(req.Form, i); err != nil {
			return err
		}

		form := req.MultipartForm.File
		if len(form) == 0 {
			return nil
		}

		ri := reflect.Indirect(reflect.ValueOf(i))
		rt := ri.Type()
		for n := range form {
			f := ri.FieldByName(n)
			if !f.IsValid() {
				for i := 0; i < rt.NumField(); i++ {
					sf := rt.Field(i)
					if sf.Tag.Get("form") == n {
						f = ri.Field(i)
						break
					}
				}
			}
			if !f.IsValid() {
				continue
			}
			if _, ok := f.Interface().(File); !ok {
				continue
			}
			mf, mh, err := req.FormFile(n)
			if err != nil {
				return err
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
