package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Password handles password hashing
type Password struct {
	Text *string // Exported so other packages can access it
	Hash []byte
}

// Set hashes a plain text password
func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.Hash = hash
	p.Text = &text
	return nil
}

// User represents a user in the system
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  Password  `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HashPassword hashes the user's password
func (u *User) HashPassword(plaintext string) error {
	return u.Password.Set(plaintext)
}

// UserPayload is the request body for user creation/login
type UserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
