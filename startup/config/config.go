package config

import "os"

type Config struct {
	Port             string
	ConnectionDBHost string
	ConnectionDBPort string
	SigningJwtKey    string
}

func NewConfig() *Config {
	return &Config{
		Port:             os.Getenv("CONNECTION_SERVICE_PORT"),
		ConnectionDBHost: os.Getenv("CONNECTION_DB_HOST"),
		ConnectionDBPort: os.Getenv("CONNECTION_DB_PORT"),
		SigningJwtKey:    os.Getenv("SIGNING_JWT_KEY"),
	}
}
