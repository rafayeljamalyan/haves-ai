package main

import (
	"log"

	"haves/internal/server"
)

func main() {
	srv := server.New()

	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
