package infrastructure_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/configs"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/infrastructure"
	"github.com/stretchr/testify/require"
)

func setupPVZTestDB(t *testing.T) *sql.DB {
	dsn := configs.GetTestPGDSN()
	if dsn == "" {
		t.Skip("TEST_PG_DSN not set")
	}
	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS pvz (
    id UUID PRIMARY KEY,
    registration_date TIMESTAMPTZ NOT NULL,
    city TEXT NOT NULL
);
DELETE FROM pvz;
`)
	require.NoError(t, err)
	return db
}

func TestPGPVZRepository_SaveAndList(t *testing.T) {
	db := setupPVZTestDB(t)
	repo := infrastructure.NewPGPVZRepository(db)
	ctx := context.Background()

	pvz1 := entities.PVZ{
		ID:               uuid.New(),
		RegistrationDate: time.Now().Add(-24 * time.Hour).UTC(),
		City:             "Москва",
	}
	pvz2 := entities.PVZ{
		ID:               uuid.New(),
		RegistrationDate: time.Now().UTC(),
		City:             "Казань",
	}

	// Act: save
	_, err := repo.Save(ctx, pvz1)
	require.NoError(t, err)
	_, err = repo.Save(ctx, pvz2)
	require.NoError(t, err)

	// Act: list all
	res, err := repo.List(ctx, nil, nil, 1, 10)
	require.NoError(t, err)
	require.Len(t, res, 2)

	// Act: filter by date
	start := time.Now().Add(-2 * time.Hour)
	res, err = repo.List(ctx, &start, nil, 1, 10)
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Equal(t, pvz2.ID, res[0].ID)

	// Act: pagination
	res, err = repo.List(ctx, nil, nil, 2, 1)
	require.NoError(t, err)
	require.Len(t, res, 1)
}
