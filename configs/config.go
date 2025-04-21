package configs

import (
	"os"
)

// Config — конфиг приложения
type Config struct {
	JWTSecret string
}

// LoadConfig загружает конфиг из переменных окружения
func LoadConfig() *Config {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET env var is required")
	}
	return &Config{
		JWTSecret: secret,
	}
}
