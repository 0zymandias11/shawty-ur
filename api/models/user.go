package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Password handles password hashing
type Password struct {
	Text string // Exported so other packages can access it
	Hash []byte
}

// Set hashes a plain text password
// func (p *Password) Set(text string) error {

// }

// User represents a user in the system
type User struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Password      Password  `json:"-"`                    // Don't expose password in JSON
	Provider      string    `json:"provider"`             // 'local' or 'google'
	ProviderID    *string   `json:"provider_id,omitempty"` // Google's user ID (nullable)
	AvatarURL     *string   `json:"avatar_url,omitempty"`  // Profile picture URL (nullable)
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

// HashPassword hashes the user's password
func (u *User) HashPassword(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password.Hash = hash
	u.Password.Text = plaintext // Store for later comparison
	return nil
}

func (u *User) CompareHash() error {
	err := bcrypt.CompareHashAndPassword(u.Password.Hash, []byte(u.Password.Text))
	if err != nil {
		return err
	}
	return nil
}

// UserPayload is the request body for user creation/login
type UserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
