// config/config.go
package config

import (
	"os"
)

type Config struct {
	DBType      string
	DBUser      string
	DBPassword  string
	DBHost      string
	DBPort      string
	DBName      string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	JWTSecret string
}

func LoadConfig() Config {
	return Config{
		DBType:      os.Getenv("DB_TYPE"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBName:      os.Getenv("DB_NAME"),

		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       0,

		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
