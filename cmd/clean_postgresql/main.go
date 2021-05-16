package main

import (
	"fmt"
	"github.com/Viverov/guideliner/internal/config"
	"github.com/Viverov/guideliner/internal/db"
	"os"
	"strings"
)

func main() {
	env := strings.ToUpper(os.Getenv("GUIDELINER_ENV"))
	cfg := config.InitConfig(env, "./config.json")

	if env == "" {
		panic("GUIDELINER_ENV must be defined for clean db")
	}

	if env == config.EnvProduction {
		panic("Can't clean DB in production mode")
	}

	fmt.Printf("Clean database %s...\n", cfg.DB.Name)
	dbInstance := db.GetDB(&db.DBOptions{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Login:    cfg.DB.Login,
		Password: cfg.DB.Password,
		Name:     cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
	})

	db, err := dbInstance.DB()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Drop database %s...\n", cfg.DB.Name)
	_, err = db.Exec("DROP SCHEMA public CASCADE;")
	if err != nil {
		panic(err)
	}
	fmt.Println("Done!")
	fmt.Printf("Recreate database %s...\n", cfg.DB.Name)
	_, err = db.Exec("CREATE SCHEMA public;")
	if err != nil {
		panic(err)
	}
	fmt.Println("Done!")
}
