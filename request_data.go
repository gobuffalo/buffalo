package buffalo

import "sync"

type requestData struct {
	d    map[string]any
	moot *sync.RWMutex
}

func newRequestData() *requestData {
	return &requestData{
		d:    make(map[string]any),
		moot: &sync.RWMutex{},
	}
}
