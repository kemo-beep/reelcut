// Run database migrations (up) and exit. Use same env as API (e.g. .env or DATABASE_URL).
package main

import (
	"log"
	"os"

	"reelcut/internal/config"
	"reelcut/pkg/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	if cfg.Database.URL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	if err := database.RunMigrations(cfg.Database.URL); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	log.Println("migrations applied successfully")
	os.Exit(0)
}
