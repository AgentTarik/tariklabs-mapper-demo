package main

import (
	"log"

	"go-test/go-test/internal/app"
)

func main() {
	application := app.New()

	log.Println("Starting server on :8080")
	if err := application.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
