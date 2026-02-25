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

func main() {
	router := mux.NewRouter()

	cfg := db.NewConfigFromEnv()
	postgresClient := db.NewClient(cfg)
	if err := postgresClient.Init(); err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewPostgresUserRepository(postgresClient)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	userHandler.RegisterRoutes(router)

	addr := ":8081"
	log.Printf("users service listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
