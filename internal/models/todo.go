package models

import "time"

// Todo represents a task owned by a user.
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTodoInput captures payload for creating a todo.
type CreateTodoInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
}

// UpdateTodoInput captures payload for updating a todo.
type UpdateTodoInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   *bool  `json:"completed"`
}
