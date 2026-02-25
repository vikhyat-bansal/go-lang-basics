package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"go-lang-basics/internal/db"
	"go-lang-basics/internal/handlers"
	"go-lang-basics/internal/repository"
	"go-lang-basics/internal/services"
)

func newServer() (*http.Server, error) {
	router := mux.NewRouter()

	cfg := db.NewConfigFromEnv()
	postgresClient := db.NewClient(cfg)
	if err := postgresClient.Init(); err != nil {
		return nil, err
	}

	userRepo := repository.NewPostgresUserRepository(postgresClient)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	userHandler.RegisterRoutes(router)

	todoRepo := repository.NewPostgresTodoRepository(postgresClient)
	todoService := services.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)
	todoHandler.RegisterRoutes(router)

	router.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	server := &http.Server{
		Addr:    ":8080",
		Handler: loggingMiddleware(router),
	}

	return server, nil
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
