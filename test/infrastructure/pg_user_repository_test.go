package infrastructure_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/infrastructure"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func setupTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("TEST_PG_DSN")
	if dsn == "" {
		t.Skip("TEST_PG_DSN not set")
	}
	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL,
    registration_date TIMESTAMPTZ NOT NULL,
    password_hash TEXT NOT NULL
);
DELETE FROM users;
`)
	require.NoError(t, err)
	return db
}

func TestPGUserRepository_CreateAndGetByEmail(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := infrastructure.NewPGUserRepository(db)
	ctx := context.Background()
	user := entities.User{
		ID:               uuid.New(),
		Email:            "test@avito.ru",
		Role:             entities.UserRoleModerator,
		RegistrationDate: time.Now().UTC(),
	}
	password := "password123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Act: create
	saved, err := repo.Create(ctx, user, string(hash))

	// Assert
	require.NoError(t, err)
	require.Equal(t, user.Email, saved.Email)
	require.Equal(t, user.Role, saved.Role)

	// Act: get by email
	got, gotHash, err := repo.GetByEmail(ctx, user.Email)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, user.Email, got.Email)
	require.Equal(t, user.Role, got.Role)
	require.Equal(t, string(hash), gotHash)

	// Act: get by non-existent email
	none, noneHash, err := repo.GetByEmail(ctx, "notfound@avito.ru")

	// Assert
	require.NoError(t, err)
	require.Nil(t, none)
	require.Empty(t, noneHash)
}
