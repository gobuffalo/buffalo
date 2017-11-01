package ratelimiter

import "sync"

// implement Counter interface
var _ Counter = &MemoryCounter{}

// InMemory is a middleware that uses the default options with the MemoryCounter
var InMemory = Middleware(NewMemoryCounter(), &DefaultOptions)

// MemoryCounter limits calls to the api based on IP
type MemoryCounter struct {
	m       *sync.RWMutex
	counter map[string]int
}

// NewMemoryCounter creates a new MemoryCounter with default values
func NewMemoryCounter() *MemoryCounter {
	return &MemoryCounter{
		m:       &sync.RWMutex{},
		counter: map[string]int{},
	}
}

// Increment increments the counter for the IP
func (r *MemoryCounter) Increment(ip string) (count int, err error) {
	r.m.Lock()
	defer r.m.Unlock()

	r.counter[ip]++
	return r.counter[ip], nil
}

// Decrement decrements the counter for the IP
func (r *MemoryCounter) Decrement(ip string) (count int, err error) {
	r.m.Lock()
	defer r.m.Unlock()

	r.counter[ip]--
	return r.counter[ip], nil
}

// Count returns the current count
func (r *MemoryCounter) Count(ip string) (count int, err error) {
	r.m.RLock()
	defer r.m.RUnlock()

	return r.counter[ip], nil
}

// Set sets the ip count
func (r *MemoryCounter) Set(ip string, count int) error {
	r.m.Lock()
	defer r.m.Unlock()

	r.counter[ip] = count
	return nil
}
