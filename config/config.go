package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort         string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPass          string
	DBName          string
	DBSSL           string
	JWTSecret       string
	CACHE_TTL       time.Duration
	CACHE_URL       string
	CACHE_PASSWORD  string
	CACHE_DATABASE  string
	PRODUCT_API_URL string
}

func LoadConfig() *Config {
	viper.AutomaticEnv()

	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_SSL", "false")
	viper.SetDefault("CACHE_TTL", 3)
	viper.SetDefault("CACHE_URL", "redis://redis:6379")
	viper.SetDefault("CACHE_PASSWORD", "")
	viper.SetDefault("CACHE_DATABASE", "0")

	return &Config{
		AppPort:         getEnv("APP_PORT"),
		DBHost:          getEnv("DB_HOST"),
		DBPort:          getEnv("DB_PORT"),
		DBUser:          getEnv("DB_USER"),
		DBPass:          getEnv("DB_PASSWORD"),
		DBName:          getEnv("DB_NAME"),
		DBSSL:           getEnv("DB_SSL"),
		JWTSecret:       getEnv("JWT_SECRET"),
		CACHE_TTL:       time.Duration(viper.GetInt("CACHE_TTL")) * time.Minute,
		CACHE_URL:       getEnv("CACHE_URL"),
		CACHE_PASSWORD:  getEnv("CACHE_PASSWORD"),
		CACHE_DATABASE:  getEnv("CACHE_DATABASE"),
		PRODUCT_API_URL: getEnv("PRODUCT_API_URL"),
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("required environment variable not set: %s", key)
	}
	return val
}
