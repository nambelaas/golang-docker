package main

import (
	"log"

	"github.com/nambelaas/golang-docker/configs"
	appinit "github.com/nambelaas/golang-docker/internal/init"
)

func main() {
	// Load configuration
	cfg := configs.New()
	log.Println("Configuration loaded")

	// Initialize database
	db, err := appinit.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
}