package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"go-lang-basics/internal/models"
	"go-lang-basics/internal/repository"
	"go-lang-basics/internal/services"
	"go-lang-basics/internal/utils"
)

// UserHandler handles HTTP requests for users.
type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/users", h.list).Methods(http.MethodGet)
	router.HandleFunc("/users", h.create).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}", h.getByID).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", h.update).Methods(http.MethodPut)
	router.HandleFunc("/users/{id}", h.delete).Methods(http.MethodDelete)
}

func (h *UserHandler) create(w http.ResponseWriter, r *http.Request) {
	var input models.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.Create(input)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	utils.JSON(w, http.StatusCreated, user)
}

func (h *UserHandler) list(w http.ResponseWriter, _ *http.Request) {
	users, err := h.service.List()
	if err != nil {
		h.handleDomainError(w, err)
		return
	}
	utils.JSON(w, http.StatusOK, users)
}

func (h *UserHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id, ok := userIDFromPath(w, r)
	if !ok {
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	utils.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := userIDFromPath(w, r)
	if !ok {
		return
	}

	var input models.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)

	user, err := h.service.Update(id, input)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	utils.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := userIDFromPath(w, r)
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		h.handleDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) handleDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		utils.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, services.ErrInvalidName), errors.Is(err, services.ErrInvalidEmail):
		utils.Error(w, http.StatusBadRequest, err.Error())
	default:
		utils.Error(w, http.StatusInternalServerError, "internal server error")
	}
}

func userIDFromPath(w http.ResponseWriter, r *http.Request) (int, bool) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid user id")
		return 0, false
	}
	return id, true
}
