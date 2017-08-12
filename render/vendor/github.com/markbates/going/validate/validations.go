package validate

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/markbates/going/wait"
)

func init() {
	fmt.Println("This package has been deprecated. Please use github.com/markbates/validate instead.")
}

// Errors holds onto all of the error messages
// that get generated during the validation process.
type Errors struct {
	Errors map[string][]string `json:"errors"`
	Lock   *sync.RWMutex       `json:"-"`
}

// Validator must be implemented in order to pass the
// validator object into the Validate function.
type Validator interface {
	IsValid(errors *Errors)
}

type vfWrapper struct {
	vf func(errors *Errors)
}

func (v vfWrapper) IsValid(errors *Errors) {
	v.vf(errors)
}

func ValidatorFunc(fn func(errors *Errors)) Validator {
	return vfWrapper{fn}
}

// NewErrors returns a pointer to a Errors
// object that has been primed and ready to go.
func NewErrors() *Errors {
	return &Errors{
		Errors: make(map[string][]string),
		Lock:   new(sync.RWMutex),
	}
}

// Error implements the error interface
func (v *Errors) Error() string {
	errs := []string{}
	for _, v := range v.Errors {
		errs = append(errs, v...)
	}
	return strings.Join(errs, "\n")
}

// Count returns the number of errors.
func (v *Errors) Count() int {
	return len(v.Errors)
}

// HasAny returns true/false depending on whether any errors
// have been tracked.
func (v *Errors) HasAny() bool {
	return v.Count() > 0
}

// Append concatenates two Errors objects together.
// This will modify the first object in place.
func (v *Errors) Append(ers *Errors) {
	for key, value := range ers.Errors {
		for _, msg := range value {
			v.Add(key, msg)
		}
	}
}

// Add will add a new message to the list of errors using
// the given key. If the key already exists the message will
// be appended to the array of the existing messages.
func (v *Errors) Add(key string, msg string) {
	v.Lock.Lock()
	v.Errors[key] = append(v.Errors[key], msg)
	v.Lock.Unlock()
}

// Get returns an array of error messages for the given key.
func (v *Errors) Get(key string) []string {
	return v.Errors[key]
}

func (v *Errors) String() string {
	b, _ := json.Marshal(v)
	return string(b)
}

// Validate takes in n number of Validator objects and will run
// them and return back a point to a Errors object that
// will contain any errors.
func Validate(validators ...Validator) *Errors {
	errors := NewErrors()

	wait.Wait(len(validators), func(index int) {
		validator := validators[index]
		validator.IsValid(errors)
	})

	return errors
}
