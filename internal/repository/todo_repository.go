package repository

import (
	"encoding/json"
	"errors"
	"fmt"

	"go-lang-basics/internal/db"
	"go-lang-basics/internal/models"
)

var ErrTodoNotFound = errors.New("todo not found")

// TodoRepository provides persistence operations for todos.
type TodoRepository interface {
	Create(input models.CreateTodoInput) (models.Todo, error)
	List() ([]models.Todo, error)
	GetByID(id int) (models.Todo, error)
	Update(id int, input models.UpdateTodoInput) (models.Todo, error)
	Delete(id int) error
}

// PostgresTodoRepository stores todos in PostgreSQL.
type PostgresTodoRepository struct {
	client *db.Client
}

func NewPostgresTodoRepository(client *db.Client) *PostgresTodoRepository {
	return &PostgresTodoRepository{client: client}
}

func (r *PostgresTodoRepository) Create(input models.CreateTodoInput) (models.Todo, error) {
	query := fmt.Sprintf(`
		SELECT row_to_json(t)
		FROM (
			INSERT INTO todos (title, description, user_id)
			VALUES ('%s', '%s', %d)
			RETURNING id, title, description, completed, user_id, created_at, updated_at
		) t;`, db.EscapeLiteral(input.Title), db.EscapeLiteral(input.Description), input.UserID)

	payload, err := r.client.QueryValue(query)
	if err != nil {
		return models.Todo{}, err
	}

	var todo models.Todo
	if err := json.Unmarshal([]byte(payload), &todo); err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func (r *PostgresTodoRepository) List() ([]models.Todo, error) {
	query := `
		SELECT COALESCE(json_agg(t), '[]'::json)
		FROM (
			SELECT id, title, description, completed, user_id, created_at, updated_at
			FROM todos
			ORDER BY id
		) t;`

	payload, err := r.client.QueryValue(query)
	if err != nil {
		return nil, err
	}

	todos := make([]models.Todo, 0)
	if err := json.Unmarshal([]byte(payload), &todos); err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *PostgresTodoRepository) GetByID(id int) (models.Todo, error) {
	query := fmt.Sprintf(`
		SELECT COALESCE(row_to_json(t)::text, '')
		FROM (
			SELECT id, title, description, completed, user_id, created_at, updated_at
			FROM todos
			WHERE id = %d
		) t;`, id)

	payload, err := r.client.QueryValue(query)
	if err != nil {
		return models.Todo{}, err
	}
	if payload == "" {
		return models.Todo{}, ErrTodoNotFound
	}

	var todo models.Todo
	if err := json.Unmarshal([]byte(payload), &todo); err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func (r *PostgresTodoRepository) Update(id int, input models.UpdateTodoInput) (models.Todo, error) {
	completedSQL := "completed"
	if input.Completed != nil {
		if *input.Completed {
			completedSQL = "TRUE"
		} else {
			completedSQL = "FALSE"
		}
	}

	query := fmt.Sprintf(`
		SELECT COALESCE(row_to_json(t)::text, '')
		FROM (
			UPDATE todos
			SET title = COALESCE(NULLIF('%s', ''), title),
				description = COALESCE(NULLIF('%s', ''), description),
				completed = %s,
				updated_at = NOW()
			WHERE id = %d
			RETURNING id, title, description, completed, user_id, created_at, updated_at
		) t;`, db.EscapeLiteral(input.Title), db.EscapeLiteral(input.Description), completedSQL, id)

	payload, err := r.client.QueryValue(query)
	if err != nil {
		return models.Todo{}, err
	}
	if payload == "" {
		return models.Todo{}, ErrTodoNotFound
	}

	var todo models.Todo
	if err := json.Unmarshal([]byte(payload), &todo); err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func (r *PostgresTodoRepository) Delete(id int) error {
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM (
			DELETE FROM todos
			WHERE id = %d
			RETURNING id
		) t;`, id)

	result, err := r.client.QueryValue(query)
	if err != nil {
		return err
	}
	if result == "0" {
		return ErrTodoNotFound
	}
	return nil
}
