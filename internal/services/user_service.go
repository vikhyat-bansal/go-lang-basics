package services

import (
	"errors"
	"strings"

	"go-lang-basics/internal/models"
	"go-lang-basics/internal/repository"
)

var (
	ErrInvalidName  = errors.New("name is required")
	ErrInvalidEmail = errors.New("email is required")
)

// UserService contains business logic for users.
type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(input models.CreateUserInput) (models.User, error) {
	if strings.TrimSpace(input.Name) == "" {
		return models.User{}, ErrInvalidName
	}
	if strings.TrimSpace(input.Email) == "" {
		return models.User{}, ErrInvalidEmail
	}
	return s.repo.Create(input)
}

func (s *UserService) List() ([]models.User, error) {
	return s.repo.List()
}

func (s *UserService) GetByID(id int) (models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) Update(id int, input models.UpdateUserInput) (models.User, error) {
	if input.Name != "" && strings.TrimSpace(input.Name) == "" {
		return models.User{}, ErrInvalidName
	}
	if input.Email != "" && strings.TrimSpace(input.Email) == "" {
		return models.User{}, ErrInvalidEmail
	}
	return s.repo.Update(id, input)
}

func (s *UserService) Delete(id int) error {
	return s.repo.Delete(id)
}
