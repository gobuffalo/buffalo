package events

import (
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// Manager can be implemented to replace the default
// events manager
type Manager interface {
	Listen(string, Listener)
	Emit(Event) error
	StopListening(string)
}

var boss Manager = &manager{
	moot:      &sync.RWMutex{},
	listeners: map[string]Listener{},
}

type manager struct {
	moot      *sync.RWMutex
	listeners map[string]Listener
}

func (m *manager) Reset() {
	m.moot.Lock()
	m.listeners = map[string]Listener{}
	m.moot.Unlock()
}

func (m *manager) Listen(name string, l Listener) {
	m.moot.Lock()
	m.listeners[name] = l
	m.moot.Unlock()
}

func (m *manager) Emit(e Event) error {
	if err := e.Validate(); err != nil {
		return errors.WithStack(err)
	}
	m.moot.RLock()
	defer m.moot.RUnlock()
	e.Kind = strings.ToLower(e.Kind)
	for _, l := range m.listeners {
		go l(e)
	}
	return nil
}

func (m *manager) StopListening(name string) {
	m.moot.Lock()
	delete(m.listeners, name)
	m.moot.Unlock()
}
