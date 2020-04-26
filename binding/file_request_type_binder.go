package binding

import (
	"net/http"
	"reflect"

	"github.com/monoculum/formam"
)

// FileRequestTypeBinder is in charge of binding File request types.
type FileRequestTypeBinder struct{}

// RegisterTo register this RequestTypeBinder to the passed request binder
// on the File content types (multipart/form-data).
func (ht FileRequestTypeBinder) RegisterTo(binder *RequestBinder) {
	binder.Register("multipart/form-data", ht.binder(binder.formDecoder))
}

func (ht FileRequestTypeBinder) binder(decoder *formam.Decoder) Binder {
	return func(req *http.Request, i interface{}) error {
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
}
