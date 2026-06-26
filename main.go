package main

import (
	"context"
	"log"

	"haves/internal/db"
	"haves/internal/server"
)

func main() {
	ctx := context.Background()

	database, err := db.New(ctx, db.DSNFromEnv())
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer database.Close()

	srv := server.New(database)

	port := "8080"

	if err := srv.Run(":" + port); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
