package binding

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gobuffalo/buffalo/internal/httpx"
)

var (
	errBlankContentType = errors.New("blank content type")
)

// RequestBinder is in charge of binding multiple requests types to
// structs.
type RequestBinder struct {
	lock    *sync.RWMutex
	binders map[string]Binder
}

// Register maps a request Content-Type (application/json)
// to a Binder.
func (rb RequestBinder) Register(contentType string, fn Binder) {
	rb.lock.Lock()
	defer rb.lock.Unlock()

	rb.binders[strings.ToLower(contentType)] = fn
}

// Exec binds a request with a passed value, depending on the content type
// It will look for the correct RequestTypeBinder and use it.
func (rb RequestBinder) Exec(req *http.Request, value interface{}) error {
	if ba, ok := value.(Bindable); ok {
		return ba.Bind(req)
	}

	ct := httpx.ContentType(req)
	if ct == "" {
		return errBlankContentType
	}

	if b, ok := rb.binders[ct]; ok {
		return b(req, value)
	}

	return fmt.Errorf("could not find a binder for %s", ct)
}

// NewRequestBinder creates our request binder with support for
// XML, JSON, HTTP and File request types.
func NewRequestBinder(requestBinders ...ContenTypeBinder) *RequestBinder {
	result := &RequestBinder{
		lock:    &sync.RWMutex{},
		binders: map[string]Binder{},
	}

	for _, requestBinder := range requestBinders {
		contentTypes := requestBinder.ContentTypes()

		for _, contentType := range contentTypes {
			result.Register(contentType, requestBinder.BinderFunc())
		}
	}

	return result
}
