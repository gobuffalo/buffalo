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
// struct.
type RequestBinder struct {
	lock    *sync.RWMutex
	binders map[string]Binder
}

// Register maps a request Content-Type (application/json)
// to a Binder.
func (rb *RequestBinder) Register(contentType string, fn Binder) {
	rb.lock.Lock()
	defer rb.lock.Unlock()

	rb.binders[strings.ToLower(contentType)] = fn
}

// Exec binds a request with a passed value, depending on the content type
// It will look for the correct RequestTypeBinder and use it.
func (rb *RequestBinder) Exec(req *http.Request, value interface{}) error {
	rb.lock.Lock()
	defer rb.lock.Unlock()

	if ba, ok := value.(Bindable); ok {
		return ba.Bind(req)
	}

	ct := httpx.ContentType(req)
	if ct == "" {
		return errBlankContentType
	}

	binder := rb.binders[ct]
	if binder == nil {
		return fmt.Errorf("could not find a binder for %s", ct)
	}

	return binder(req, value)
}

// NewRequestBinder creates our request binder with support for
// XML, JSON, HTTP and File request types.
func NewRequestBinder(requestBinders ...ContenTypeBinder) *RequestBinder {
	result := &RequestBinder{
		lock:    &sync.RWMutex{},
		binders: map[string]Binder{},
	}

	for _, requestBinder := range requestBinders {
		for _, contentType := range requestBinder.ContentTypes() {
			result.Register(contentType, requestBinder.BinderFunc())
		}
	}

	return result
}
