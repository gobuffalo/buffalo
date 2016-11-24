package models

import (
	"encoding/json"
	"time"
)

type User struct {
	ID        int       `json:"id" db:"id"`
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

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	b, _ := json.Marshal(u)
	return string(b)
}
