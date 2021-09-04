package main

import (
	"fmt"
	"github.com/Viverov/guideliner/internal/config"
	"github.com/Viverov/guideliner/internal/cradle"
	"github.com/Viverov/guideliner/internal/db"
	"github.com/Viverov/guideliner/internal/domains/guide"
	"github.com/Viverov/guideliner/internal/domains/user"
	"github.com/Viverov/guideliner/internal/server"
	"os"
	"strings"
)

func main() {
	// Setup config
	env := strings.ToUpper(os.Getenv("GUIDELINER_ENV"))
	cfg := config.InitConfig(env, "./config.json")

	// Setup cradle
	fmt.Println("Setup cradle...")
	cradleBuilder := cradle.Builder{}

	// Init config
	cradleBuilder.SetConfig(cfg)

	// Init DB
	dbInstance := db.GetDB(&db.DBOptions{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Login:    cfg.DB.Login,
		Password: cfg.DB.Password,
		Name:     cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
	})
	cradleBuilder.SetSqlDB(dbInstance)

	// Init services
	cradleBuilder.SetUserService(user.BuildUserService(cfg.Tokens.SecretKey, dbInstance))
	cradleBuilder.SetGuideService(guide.BuildGuideService(dbInstance))

	cradleObj := cradleBuilder.Build()
	fmt.Println("Done!")

	s := server.Init(cradleObj)
	s.Run()
}
