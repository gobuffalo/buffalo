package binding

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RequestBinder_Exec(t *testing.T) {
	r := require.New(t)

	var used bool
	defaultRequestBinder.Register("paganotoni/test", func(*http.Request, interface{}) error {
		used = true
		return nil
	})

	req, err := http.NewRequest("GET", "/home", strings.NewReader(""))
	req.Header.Add("content-type", "paganotoni/test")
	r.NoError(err)

	data := &struct{}{}
	r.NoError(defaultRequestBinder.Exec(req, data))
	r.True(used)
}

func Test_RequestBinder_Exec_BlankContentType(t *testing.T) {
	r := require.New(t)

	req, err := http.NewRequest("GET", "/home", strings.NewReader(""))
	r.NoError(err)

	data := &struct{}{}
	r.Equal(defaultRequestBinder.Exec(req, data), errBlankContentType)
}

func Test_RequestBinder_Exec_Bindable(t *testing.T) {
	r := require.New(t)

	defaultRequestBinder.Register("paganotoni/orbison", func(req *http.Request, val interface{}) error {
		switch v := val.(type) {
		case orbison:
			v.bound = false
		}

		return errors.New("this should not be called")
	})

	req, err := http.NewRequest("GET", "/home", strings.NewReader(""))
	req.Header.Add("content-type", "paganotoni/orbison")
	r.NoError(err)

	data := &orbison{}
	r.NoError(defaultRequestBinder.Exec(req, data))
	r.True(data.bound)
}

func Test_RequestBinder_Exec_NoBinder(t *testing.T) {
	r := require.New(t)

	req, err := http.NewRequest("GET", "/home", strings.NewReader(""))
	req.Header.Add("content-type", "paganotoni/other")
	r.NoError(err)

	err = defaultRequestBinder.Exec(req, &struct{}{})
	r.Error(err)
	r.Equal(err.Error(), "could not find a binder for paganotoni/other")
}
