package main

import (
	"github.com/Viverov/guideliner/internal/config"
	"github.com/Viverov/guideliner/internal/db"
	"github.com/Viverov/guideliner/internal/migrations"
	"github.com/go-gormigrate/gormigrate/v2"
	"log"
	"os"
	"strings"
)

func main() {
	// Setup config
	env := strings.ToUpper(os.Getenv("GUIDELINER_ENV"))
	cfg := config.InitConfig(env, "./config.json")

	dbInstance := db.GetDB(&db.DBOptions{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Login:    cfg.DB.Login,
		Password: cfg.DB.Password,
		Name:     cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
	})

	options := &gormigrate.Options{
		UseTransaction:            false,
		ValidateUnknownMigrations: true,
	}

	migrator := gormigrate.New(dbInstance, options, migrations.GetMigrationsList())
	err := migrator.Migrate()
	if err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Printf("Migration did run successfully")
}
