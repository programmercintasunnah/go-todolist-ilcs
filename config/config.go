// config/config.go
package config

import (
	"os"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBService  string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	JWTSecret string
}

func LoadConfig() Config {
	return Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBService:  os.Getenv("DB_SERVICE"),

		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       0, // Default DB

		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
