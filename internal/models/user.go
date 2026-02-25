package models

import "time"

// User represents a basic user record.
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserInput captures payload for creating a user.
type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateUserInput captures payload for updating a user.
type UpdateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
