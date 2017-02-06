package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
)

// User model stores information about a user account
type User struct {
	ID        int       `json:"id" db:"id" schema:"-"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Email     string    `json:"email" db:"email"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	b, _ := json.Marshal(u)
	return string(b)
}

// FullName returns the first name and last name, separated by a space
func (u User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// ValidateNew validates the User fields and checks whether a user has
// already claimed that email address
func (u *User) ValidateNew(tx *pop.Connection) (*validate.Errors, error) {
	verrs, err := u.validateCommon(tx)
	verrs.Append(validate.Validate(
		&validators.FuncValidator{
			Fn: func() bool {
				var b bool
				if u.Email != "" {
					b, err = tx.Where("email = ?", u.Email).Exists(u)
				}
				return !b
			},
			Field:   "Email",
			Message: "%s was already taken.",
		},
	))
	return verrs, err
}

// ValidateUpdate validates the User fields and confirms that the email has
// not been claimed by a different user
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	verrs, err := u.validateCommon(tx)
	verrs.Append(validate.Validate(
		&validators.FuncValidator{
			Fn: func() bool {
				var b bool
				if u.Email != "" {
					// Differs from ValidateNew here - only check whether *other* users claimed email
					b, err = tx.Where("email = ? and id != ?", u.Email, u.ID).Exists(u)
				}
				return !b
			},
			Field:   "Email",
			Message: "%s was already taken.",
		},
	))
	return verrs, err
}

func (u *User) validateCommon(tx *pop.Connection) (*validate.Errors, error) {
	verrs := validate.Validate(
		&validators.StringIsPresent{Name: "First Name", Field: u.FirstName},
		&validators.StringIsPresent{Name: "Last Name", Field: u.LastName},
		&validators.StringIsPresent{Name: "Email", Field: u.Email},
	)
	return verrs, nil
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	b, _ := json.Marshal(u)
	return string(b)
}
