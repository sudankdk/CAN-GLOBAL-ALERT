package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	Redis  RedisConfig
}

type ServerConfig struct {
	Port string
}

type RedisConfig struct {
	Host string
	Port string
}

func Load() *Config {
	// Load .env only if present (dev-friendly, prod-safe)
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		Redis: RedisConfig{
			Host: getEnv("REDIS_HOST", "localhost"),
			Port: getEnv("REDIS_PORT", "6379"),
		},
	}

	validate(cfg)
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func validate(cfg *Config) {
	if cfg.Server.Port == "" {
		log.Fatal("server port cannot be empty")
	}
	if cfg.Redis.Host == "" {
		log.Fatal("redis host cannot be empty")
	}
	if cfg.Redis.Port == "" {
		log.Fatal("redis port cannot be empty")
	}
}
