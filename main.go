package main

import (
	"log"

	"haves/internal/server"
)

func main() {
	srv := server.New()

	port := "8081"

	if err := srv.Run(":" + port); err != nil {
		log.Fatalf("server failed to start: %v", err)
	} else {
		log.Printf("server started on port %v", port)
	}
}
