package buffalo

import "sync"

type requestData struct {
	d    map[string]interface{}
	moot *sync.RWMutex
}

func newRequestData() *requestData {
	return &requestData{
		d:    make(map[string]interface{}),
		moot: &sync.RWMutex{},
	}
}
