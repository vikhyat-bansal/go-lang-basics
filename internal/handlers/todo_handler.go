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

// TodoHandler handles HTTP requests for todos.
type TodoHandler struct {
	service *services.TodoService
}

func NewTodoHandler(service *services.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

func (h *TodoHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/todos", h.list).Methods(http.MethodGet)
	router.HandleFunc("/todos", h.create).Methods(http.MethodPost)
	router.HandleFunc("/todos/{id}", h.getByID).Methods(http.MethodGet)
	router.HandleFunc("/todos/{id}", h.update).Methods(http.MethodPut)
	router.HandleFunc("/todos/{id}", h.delete).Methods(http.MethodDelete)
}

func (h *TodoHandler) create(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTodoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	todo, err := h.service.Create(input)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}
	utils.JSON(w, http.StatusCreated, todo)
}

func (h *TodoHandler) list(w http.ResponseWriter, _ *http.Request) {
	todos, err := h.service.List()
	if err != nil {
		h.handleDomainError(w, err)
		return
	}
	utils.JSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id, ok := todoIDFromPath(w, r)
	if !ok {
		return
	}
	todo, err := h.service.GetByID(id)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}
	utils.JSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := todoIDFromPath(w, r)
	if !ok {
		return
	}

	var input models.UpdateTodoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	todo, err := h.service.Update(id, input)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}
	utils.JSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := todoIDFromPath(w, r)
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		h.handleDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TodoHandler) handleDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrTodoNotFound):
		utils.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, services.ErrInvalidTitle), errors.Is(err, services.ErrInvalidUserID):
		utils.Error(w, http.StatusBadRequest, err.Error())
	default:
		utils.Error(w, http.StatusInternalServerError, "internal server error")
	}
}

func todoIDFromPath(w http.ResponseWriter, r *http.Request) (int, bool) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid todo id")
		return 0, false
	}
	return id, true
}
