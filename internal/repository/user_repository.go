package repository

import (
	"errors"
	"sync"
	"time"

	"go-lang-basics/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

// UserRepository provides persistence operations for users.
type UserRepository interface {
	Create(input models.CreateUserInput) models.User
	List() []models.User
	GetByID(id int) (models.User, error)
	Update(id int, input models.UpdateUserInput) (models.User, error)
	Delete(id int) error
}

// InMemoryUserRepository stores users in memory.
type InMemoryUserRepository struct {
	mu     sync.RWMutex
	nextID int
	users  map[int]models.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		nextID: 1,
		users:  make(map[int]models.User),
	}
}

func (r *InMemoryUserRepository) Create(input models.CreateUserInput) models.User {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	user := models.User{
		ID:        r.nextID,
		Name:      input.Name,
		Email:     input.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	r.users[user.ID] = user
	r.nextID++
	return user
}

func (r *InMemoryUserRepository) List() []models.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]models.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users
}

func (r *InMemoryUserRepository) GetByID(id int) (models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return models.User{}, ErrUserNotFound
	}
	return user, nil
}

func (r *InMemoryUserRepository) Update(id int, input models.UpdateUserInput) (models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return models.User{}, ErrUserNotFound
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	user.UpdatedAt = time.Now().UTC()

	r.users[id] = user
	return user, nil
}

func (r *InMemoryUserRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return ErrUserNotFound
	}
	delete(r.users, id)
	return nil
}
