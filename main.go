package main

import "log"

func main() {
	server, err := newServer()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
