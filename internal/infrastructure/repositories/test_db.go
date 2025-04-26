package repositories

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewTestPostgresDB создает новое подключение к тестовой базе данных
func NewTestPostgresDB() (*sql.DB, error) {
	dsn := os.Getenv("TEST_PG_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("TEST_PG_DSN not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Создаем необходимые таблицы
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		);

		CREATE TABLE IF NOT EXISTS pvz (
			id UUID PRIMARY KEY,
			registration_date TIMESTAMPTZ NOT NULL,
			city TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS reception (
			id UUID PRIMARY KEY,
			pvz_id UUID NOT NULL REFERENCES pvz(id),
			status TEXT NOT NULL,
			date_time TIMESTAMPTZ NOT NULL
		);

		CREATE TABLE IF NOT EXISTS product (
			id UUID PRIMARY KEY,
			reception_id UUID NOT NULL REFERENCES reception(id),
			type TEXT NOT NULL,
			received_at TIMESTAMPTZ NOT NULL
		);

		-- Очищаем таблицы перед каждым тестом
		DELETE FROM product;
		DELETE FROM reception;
		DELETE FROM pvz;
		DELETE FROM users;
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create test tables: %w", err)
	}

	return db, nil
}
