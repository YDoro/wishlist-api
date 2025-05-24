package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort   string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSL     string
	JWTSecret string
}

func LoadConfig() *Config {
	viper.AutomaticEnv()

	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_SSL", "false")

	return &Config{
		AppPort:   getEnv("APP_PORT"),
		DBHost:    getEnv("DB_HOST"),
		DBPort:    getEnv("DB_PORT"),
		DBUser:    getEnv("DB_USER"),
		DBPass:    getEnv("DB_PASSWORD"),
		DBName:    getEnv("DB_NAME"),
		DBSSL:     getEnv("DB_SSL"),
		JWTSecret: getEnv("JWT_SECRET"),
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("required environment variable not set: %s", key)
	}
	return val
}
