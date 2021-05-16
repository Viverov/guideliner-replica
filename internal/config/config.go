package config

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
	"os"
	"strings"
)

const (
	EnvDebug       = "DEBUG"
	EnvTest        = "TEST"
	EnvDevelopment = "DEVELOPMENT"
	EnvProduction  = "PRODUCTION"
)

type Config struct {
	Env    string `validate:"required"`
	Server struct {
		Host string `json:"host" envconfig:"GUIDELINER_SERVER_HOST" validate:"required"`
		Port string `json:"port" envconfig:"GUIDELINER_SERVER_PORT" validate:"required"`
	} `json:"server"`
	DB struct {
		Host     string `json:"host" envconfig:"GUIDELINER_DB_HOST" validate:"required"`
		Port     string `json:"port" envconfig:"GUIDELINER_DB_PORT" validate:"required"`
		Login    string `json:"login" envconfig:"GUIDELINER_DB_LOGIN" validate:"required"`
		Password string `json:"password" envconfig:"GUIDELINER_DB_PASSWORD" validate:"required"`
		Name     string `json:"name" envconfig:"GUIDELINER_DB_NAME" validate:"required"`
		SSLMode  string `json:"sslMode" envconfig:"GUIDELINER_DB_SSL_MODE" validate:"required"`
	} `json:"db"`
	Tokens struct {
		SecretKey string `json:"secretKey" envconfig:"GUIDELINER_TOKEN_SECRET" validate:"required"`
	} `json:"tokens"`
}

func InitConfig(env string, jsonPath string) *Config {
	cfg := &Config{}

	resolvedEnv := strings.ToUpper(env)
	if resolvedEnv == "" {
		fmt.Printf("ENV not set, use %s mode\n", EnvProduction)
		resolvedEnv = EnvProduction
	} else if resolvedEnv != EnvDebug &&
		resolvedEnv != EnvTest &&
		resolvedEnv != EnvDevelopment &&
		resolvedEnv != EnvProduction {
		panic(fmt.Sprintf("Env set, but by undefined mode: %s\n", env))
	}
	fmt.Printf("Resolve %s mode\n", resolvedEnv)
	cfg.Env = resolvedEnv

	if resolvedEnv == EnvDebug || resolvedEnv == EnvTest || resolvedEnv == EnvDevelopment {
		fmt.Printf("In %s mode, application can process with json and env config\n", resolvedEnv)
		if jsonPath == "" {
			fmt.Printf("But, jsonPath not provided, skip JSON step\n")
		} else {
			fmt.Printf("Process JSON...\n")
			readFile(cfg, jsonPath)
			fmt.Printf("Done!\n")
		}
	} else {
		fmt.Printf("In %s mode, application can process only by env config\n", resolvedEnv)
	}

	fmt.Printf("Process ENV...\n")
	readEnv(cfg)
	fmt.Printf("Done!\n")

	fmt.Printf("Validate config...\n")
	err := validator.New().Struct(cfg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Done!\n")

	return cfg
}

func readFile(cfg *Config, jsonPath string) {
	f, err := os.Open(jsonPath)
	if err != nil {
		fmt.Printf("Probles occured while file opening: %s", err.Error())
		return
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		panic(err)
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(err)
	}
}
