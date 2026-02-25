package main

import (
	"log"
)

func main() {
	server := newServer()
	log.Printf("starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
