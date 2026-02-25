package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"go-lang-basics/internal/handlers"
	"go-lang-basics/internal/repository"
	"go-lang-basics/internal/services"
)

func newServer() *http.Server {
	router := mux.NewRouter()

	userRepo := repository.NewInMemoryUserRepository()
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	userHandler.RegisterRoutes(router)

	router.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	return &http.Server{
		Addr:    ":8080",
		Handler: loggingMiddleware(router),
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
