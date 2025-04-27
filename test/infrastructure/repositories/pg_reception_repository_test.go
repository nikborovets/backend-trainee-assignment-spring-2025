package infrastructure_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/configs"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/infrastructure/repositories"
	"github.com/stretchr/testify/require"
)

func init() {
	// Загружаем переменные окружения из .env файла
	_ = godotenv.Load("../../../.env")

	// Вывод значения для отладки
	dsn := os.Getenv("TEST_PG_DSN")
	if dsn != "" {
		println("TEST_PG_DSN loaded in pg_reception:", dsn)
	}
}

func setupReceptionTestDB(t *testing.T) *sql.DB {
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
    date_time TIMESTAMPTZ NOT NULL
);
-- Очищаем таблицы в правильном порядке с учетом внешних ключей
DELETE FROM product;
DELETE FROM reception;
DELETE FROM pvz;
`)
	require.NoError(t, err)
	return db
}

func TestPGReceptionRepository_Save_GetActive_CloseLast(t *testing.T) {
	// Arrange
	db := setupReceptionTestDB(t)
	repo := repositories.NewPGReceptionRepository(db)
	ctx := context.Background()
	pvzID := uuid.New()
	_, err := db.Exec(`INSERT INTO pvz (id, registration_date, city) VALUES ($1, $2, $3)`, pvzID, time.Now().UTC(), "Москва")
	require.NoError(t, err)

	rec := entities.Reception{
		ID:       uuid.New(),
		PVZID:    pvzID,
		Status:   entities.ReceptionInProgress,
		DateTime: time.Now().UTC(),
	}

	// Act: save
	saved, err := repo.Save(ctx, rec)
	require.NoError(t, err)
	require.Equal(t, rec.ID, saved.ID)

	// Act: get active
	got, err := repo.GetActive(ctx, pvzID)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, rec.ID, got.ID)
	require.Equal(t, entities.ReceptionInProgress, got.Status)

	// Act: close last
	closedAt := time.Now().Add(1 * time.Hour).UTC()
	err = repo.CloseLast(ctx, pvzID, closedAt)
	require.NoError(t, err)

	// Assert: get active returns nil
	got, err = repo.GetActive(ctx, pvzID)
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestPGReceptionRepository_ListByPVZ(t *testing.T) {
	// Arrange
	db := setupReceptionTestDB(t)
	repo := repositories.NewPGReceptionRepository(db)
	ctx := context.Background()
	pvzID := uuid.New()
	_, err := db.Exec(`INSERT INTO pvz (id, registration_date, city) VALUES ($1, $2, $3)`, pvzID, time.Now().UTC(), "Москва")
	require.NoError(t, err)

	rec1 := entities.Reception{
		ID:       uuid.New(),
		PVZID:    pvzID,
		Status:   entities.ReceptionInProgress,
		DateTime: time.Now().Add(-2 * time.Hour).UTC(),
	}
	rec2 := entities.Reception{
		ID:       uuid.New(),
		PVZID:    pvzID,
		Status:   entities.ReceptionClosed,
		DateTime: time.Now().Add(-1 * time.Hour).UTC(),
	}
	_, err = repo.Save(ctx, rec1)
	require.NoError(t, err)
	_, err = repo.Save(ctx, rec2)
	require.NoError(t, err)

	// Act
	list, err := repo.ListByPVZ(ctx, pvzID)

	// Assert
	require.NoError(t, err)
	require.Len(t, list, 2)
	ids := []uuid.UUID{list[0].ID, list[1].ID}
	require.Contains(t, ids, rec1.ID)
	require.Contains(t, ids, rec2.ID)
}
