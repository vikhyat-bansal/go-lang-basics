package services

import (
	"errors"
	"strings"

	"go-lang-basics/internal/models"
	"go-lang-basics/internal/repository"
)

var ErrInvalidTitle = errors.New("title is required")
var ErrInvalidUserID = errors.New("user_id must be greater than 0")

// TodoService contains business logic for todos.
type TodoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) Create(input models.CreateTodoInput) (models.Todo, error) {
	if strings.TrimSpace(input.Title) == "" {
		return models.Todo{}, ErrInvalidTitle
	}
	if input.UserID <= 0 {
		return models.Todo{}, ErrInvalidUserID
	}
	return s.repo.Create(input)
}

func (s *TodoService) List() ([]models.Todo, error) {
	return s.repo.List()
}

func (s *TodoService) GetByID(id int) (models.Todo, error) {
	return s.repo.GetByID(id)
}

func (s *TodoService) Update(id int, input models.UpdateTodoInput) (models.Todo, error) {
	if input.Title != "" && strings.TrimSpace(input.Title) == "" {
		return models.Todo{}, ErrInvalidTitle
	}
	return s.repo.Update(id, input)
}

func (s *TodoService) Delete(id int) error {
	return s.repo.Delete(id)
}
