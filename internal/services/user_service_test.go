package services

import (
	"errors"
	"testing"

	"go-lang-basics/internal/models"
	"go-lang-basics/internal/repository"
)

type fakeUserRepo struct {
	createFn func(input models.CreateUserInput) (models.User, error)
}

func (f *fakeUserRepo) Create(input models.CreateUserInput) (models.User, error) {
	return f.createFn(input)
}
func (f *fakeUserRepo) List() ([]models.User, error) { return nil, nil }
func (f *fakeUserRepo) GetByID(id int) (models.User, error) {
	return models.User{}, repository.ErrUserNotFound
}
func (f *fakeUserRepo) Update(id int, input models.UpdateUserInput) (models.User, error) {
	return models.User{}, repository.ErrUserNotFound
}
func (f *fakeUserRepo) Delete(id int) error { return repository.ErrUserNotFound }

func TestUserServiceCreateValidation(t *testing.T) {
	svc := NewUserService(&fakeUserRepo{createFn: func(input models.CreateUserInput) (models.User, error) {
		return models.User{ID: 1, Name: input.Name, Email: input.Email}, nil
	}})

	if _, err := svc.Create(models.CreateUserInput{Name: "", Email: "a@b.com"}); !errors.Is(err, ErrInvalidName) {
		t.Fatalf("expected ErrInvalidName, got %v", err)
	}

	if _, err := svc.Create(models.CreateUserInput{Name: "Test", Email: ""}); !errors.Is(err, ErrInvalidEmail) {
		t.Fatalf("expected ErrInvalidEmail, got %v", err)
	}
}

func TestUserServiceCreatePassesToRepo(t *testing.T) {
	svc := NewUserService(&fakeUserRepo{createFn: func(input models.CreateUserInput) (models.User, error) {
		return models.User{ID: 7, Name: input.Name, Email: input.Email}, nil
	}})

	user, err := svc.Create(models.CreateUserInput{Name: "Alice", Email: "alice@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != 7 || user.Name != "Alice" {
		t.Fatalf("unexpected user returned: %+v", user)
	}
}
