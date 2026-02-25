package repository

import (
	"encoding/json"
	"errors"
	"fmt"

	"go-lang-basics/internal/db"
	"go-lang-basics/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

// UserRepository provides persistence operations for users.
type UserRepository interface {
	Create(input models.CreateUserInput) (models.User, error)
	List() ([]models.User, error)
	GetByID(id int) (models.User, error)
	Update(id int, input models.UpdateUserInput) (models.User, error)
	Delete(id int) error
}

// PostgresUserRepository stores users in PostgreSQL.
type PostgresUserRepository struct {
	client *db.Client
}

func NewPostgresUserRepository(client *db.Client) *PostgresUserRepository {
	return &PostgresUserRepository{client: client}
}

func (r *PostgresUserRepository) Create(input models.CreateUserInput) (models.User, error) {
	query := fmt.Sprintf(`
		SELECT row_to_json(t)
		FROM (
			INSERT INTO users (name, email)
			VALUES ('%s', '%s')
			RETURNING id, name, email, created_at, updated_at
		) t;`, db.EscapeLiteral(input.Name), db.EscapeLiteral(input.Email))

	payload, err := r.client.QueryValue(query)
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	if err := json.Unmarshal([]byte(payload), &user); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *PostgresUserRepository) List() ([]models.User, error) {
	query := `
		SELECT COALESCE(json_agg(t), '[]'::json)
		FROM (
			SELECT id, name, email, created_at, updated_at
			FROM users
			ORDER BY id
		) t;`

	payload, err := r.client.QueryValue(query)
	if err != nil {
		return nil, err
	}

	users := make([]models.User, 0)
	if err := json.Unmarshal([]byte(payload), &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *PostgresUserRepository) GetByID(id int) (models.User, error) {
	query := fmt.Sprintf(`
		SELECT COALESCE(row_to_json(t)::text, '')
		FROM (
			SELECT id, name, email, created_at, updated_at
			FROM users
			WHERE id = %d
		) t;`, id)

	payload, err := r.client.QueryValue(query)
	if err != nil {
		return models.User{}, err
	}
	if payload == "" {
		return models.User{}, ErrUserNotFound
	}

	var user models.User
	if err := json.Unmarshal([]byte(payload), &user); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *PostgresUserRepository) Update(id int, input models.UpdateUserInput) (models.User, error) {
	query := fmt.Sprintf(`
		SELECT COALESCE(row_to_json(t)::text, '')
		FROM (
			UPDATE users
			SET name = COALESCE(NULLIF('%s', ''), name),
				email = COALESCE(NULLIF('%s', ''), email),
				updated_at = NOW()
			WHERE id = %d
			RETURNING id, name, email, created_at, updated_at
		) t;`, db.EscapeLiteral(input.Name), db.EscapeLiteral(input.Email), id)

	payload, err := r.client.QueryValue(query)
	if err != nil {
		return models.User{}, err
	}
	if payload == "" {
		return models.User{}, ErrUserNotFound
	}

	var user models.User
	if err := json.Unmarshal([]byte(payload), &user); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *PostgresUserRepository) Delete(id int) error {
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM (
			DELETE FROM users
			WHERE id = %d
			RETURNING id
		) t;`, id)

	result, err := r.client.QueryValue(query)
	if err != nil {
		return err
	}
	if result == "0" {
		return ErrUserNotFound
	}
	return nil
}
