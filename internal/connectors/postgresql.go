package connectors

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBOptions struct {
	Host     string
	Port     string
	Login    string
	Password string
	Name     string
	SSLMode  string
}

func GetDB(options *DBOptions) *gorm.DB {
	if options.Host == "" ||
		options.Port == "" ||
		options.Login == "" ||
		options.Password == "" ||
		options.Name == "" ||
		options.SSLMode == "" {
		panic("Some of required params in DB creation is empty. Required params: Host,Port,User,Password,Name,SSLMode('disable' or 'enable')")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		options.Host,
		options.Login,
		options.Password,
		options.Name,
		options.Port,
		options.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
