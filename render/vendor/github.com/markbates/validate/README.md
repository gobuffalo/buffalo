# github.com/markbates/validate
[![Build Status](https://travis-ci.org/markbates/validate.svg?branch=master)](https://travis-ci.org/markbates/validate)

This package provides a framework for writing validations for Go applications. It does not, however, provide you with any actual validators, that part is up to you.

## Installation

```bash
$ go get github.com/markbates/validate
```

## Usage

Using validate is pretty easy, just define some `Validator` objects and away you go.

Here is a pretty simple example:

```go
package main

import (
	"log"

	v "github.com/markbates/validate"
)

type User struct {
	Name  string
	Email string
}

func (u *User) IsValid(errors *v.Errors) {
	if u.Name == "" {
		errors.Add("name", "Name must not be blank!")
	}
	if u.Email == "" {
		errors.Add("email", "Email must not be blank!")
	}
}

func main() {
	u := User{Name: "", Email: ""}
	errors := v.Validate(&u)
	log.Println(errors.Errors)
  // map[name:[Name must not be blank!] email:[Email must not be blank!]]
}
```

In the previous example I wrote a single `Validator` for the `User` struct. To really get the benefit of using go-validator, as well as the Go language, I would recommend creating distinct validators for each thing you want to validate, that way they can be run concurrently.

```go
package main

import (
	"fmt"
	"log"
	"strings"

	v "github.com/markbates/validate"
)

type User struct {
	Name  string
	Email string
}

type PresenceValidator struct {
	Field string
	Value string
}

func (v *PresenceValidator) IsValid(errors *v.Errors) {
	if v.Value == "" {
		errors.Add(strings.ToLower(v.Field), fmt.Sprintf("%s must not be blank!", v.Field))
	}
}

func main() {
	u := User{Name: "", Email: ""}
	errors := v.Validate(&PresenceValidator{"Email", u.Email}, &PresenceValidator{"Name", u.Name})
	log.Println(errors.Errors)
  // map[name:[Name must not be blank!] email:[Email must not be blank!]]
}
```

That's really it. Pretty simple and straight-forward Just a nice clean framework for writing your own validators. Use in good health.
