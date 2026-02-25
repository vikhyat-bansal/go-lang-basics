package services

import (
	"errors"
	"testing"

	"go-lang-basics/internal/models"
	"go-lang-basics/internal/repository"
)

type fakeTodoRepo struct {
	createFn func(input models.CreateTodoInput) (models.Todo, error)
}

func (f *fakeTodoRepo) Create(input models.CreateTodoInput) (models.Todo, error) {
	return f.createFn(input)
}
func (f *fakeTodoRepo) List() ([]models.Todo, error) { return nil, nil }
func (f *fakeTodoRepo) GetByID(id int) (models.Todo, error) {
	return models.Todo{}, repository.ErrTodoNotFound
}
func (f *fakeTodoRepo) Update(id int, input models.UpdateTodoInput) (models.Todo, error) {
	return models.Todo{}, repository.ErrTodoNotFound
}
func (f *fakeTodoRepo) Delete(id int) error { return repository.ErrTodoNotFound }

func TestTodoServiceCreateValidation(t *testing.T) {
	svc := NewTodoService(&fakeTodoRepo{createFn: func(input models.CreateTodoInput) (models.Todo, error) {
		return models.Todo{ID: 1, Title: input.Title, UserID: input.UserID}, nil
	}})

	if _, err := svc.Create(models.CreateTodoInput{Title: "", UserID: 1}); !errors.Is(err, ErrInvalidTitle) {
		t.Fatalf("expected ErrInvalidTitle, got %v", err)
	}

	if _, err := svc.Create(models.CreateTodoInput{Title: "Task", UserID: 0}); !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}
}

func TestTodoServiceCreatePassesToRepo(t *testing.T) {
	svc := NewTodoService(&fakeTodoRepo{createFn: func(input models.CreateTodoInput) (models.Todo, error) {
		return models.Todo{ID: 5, Title: input.Title, UserID: input.UserID}, nil
	}})

	todo, err := svc.Create(models.CreateTodoInput{Title: "Write tests", UserID: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if todo.ID != 5 || todo.Title != "Write tests" {
		t.Fatalf("unexpected todo returned: %+v", todo)
	}
}
