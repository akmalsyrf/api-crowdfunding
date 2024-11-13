package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvType string

const (
	Development EnvType = "development"
	Production  EnvType = "production"
	Test        EnvType = "test"
)

func (env EnvType) isValid() bool {
	switch env {
	case Development, Production, Test:
		return true
	default:
		return false
	}
}

type Service struct {
	Port int
	Env  EnvType
}

type Database struct {
	DbHost string
	DbUser string
	DbPass string
	DbName string
	DbPort string
}

type Midtrans struct {
	MidtransServiceKey string
	MidtransClientKey  string
}

type Config struct {
	Service  Service
	Database Database
	Midtrans Midtrans
}

func LoadEnv() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	portStr := os.Getenv("PORT")
	env := EnvType(os.Getenv("ENV"))
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	midtransServiceKey := os.Getenv("MIDTRANS_SERVER_KEY")
	midtransClientKey := os.Getenv("MIDTRANS_CLIENT_KEY")

	if !env.isValid() {
		return nil, fmt.Errorf("invalid env value: %s", env)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port value: %v", err)
	}

	config := Config{
		Service: Service{
			Port: port,
			Env:  env,
		},
		Database: Database{
			DbHost: dbHost,
			DbUser: dbUser,
			DbPass: dbPass,
			DbName: dbName,
			DbPort: dbPort,
		},
		Midtrans: Midtrans{
			MidtransServiceKey: midtransServiceKey,
			MidtransClientKey:  midtransClientKey,
		},
	}

	return &config, nil
}
