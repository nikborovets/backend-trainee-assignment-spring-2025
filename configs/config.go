package configs

import (
	"os"
)

// Config — конфиг приложения
type Config struct {
	JWTSecret string
	PGDSN     string
}

// LoadConfig загружает конфиг из переменных окружения
func LoadConfig() *Config {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET env var is required")
	}
	pgDsn := os.Getenv("PG_DSN")
	if pgDsn == "" {
		panic("PG_DSN env var is required")
	}
	return &Config{
		JWTSecret: secret,
		PGDSN:     pgDsn,
	}
}

// GetTestPGDSN возвращает строку подключения к тестовой БД из переменной окружения TEST_PG_DSN
func GetTestPGDSN() string {
	return os.Getenv("TEST_PG_DSN")
}
