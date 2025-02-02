package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MySQLDSN  string
	RedisAddr string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := &Config{
		MySQLDSN:  os.Getenv("MYSQL_DSN"),
		RedisAddr: os.Getenv("REDIS_ADDR"),
	}

	if cfg.MySQLDSN == "" || cfg.RedisAddr == "" {
		log.Fatal("Missing required environment variables")
	}

	return cfg
}
