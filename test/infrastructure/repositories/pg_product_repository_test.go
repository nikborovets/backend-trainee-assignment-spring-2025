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
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/infrastructure/repositories"
	"github.com/stretchr/testify/require"
)

func setupProductTestDB(t *testing.T) *sql.DB {
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
    received_at TIMESTAMPTZ NOT NULL
);
DELETE FROM product;
DELETE FROM reception;
DELETE FROM pvz;
`)
	require.NoError(t, err)
	return db
}

func TestPGProductRepository_Save_DeleteLast_ListByReception(t *testing.T) {
	// Arrange
	db := setupProductTestDB(t)
	repo := repositories.NewPGProductRepository(db)
	ctx := context.Background()
	pvzID := uuid.New()
	_, err := db.Exec(`INSERT INTO pvz (id, registration_date, city) VALUES ($1, $2, $3)`, pvzID, time.Now().UTC(), "Москва")
	require.NoError(t, err)
	recID := uuid.New()
	_, err = db.Exec(`INSERT INTO reception (id, pvz_id, status, date_time) VALUES ($1, $2, $3, $4)`, recID, pvzID, "in_progress", time.Now().UTC())
	require.NoError(t, err)

	p1 := entities.Product{
		ID:          uuid.New(),
		ReceptionID: recID,
		Type:        entities.ProductElectronics,
		ReceivedAt:  time.Now().Add(-2 * time.Hour).UTC(),
	}
	p2 := entities.Product{
		ID:          uuid.New(),
		ReceptionID: recID,
		Type:        entities.ProductClothes,
		ReceivedAt:  time.Now().Add(-1 * time.Hour).UTC(),
	}
	p3 := entities.Product{
		ID:          uuid.New(),
		ReceptionID: recID,
		Type:        entities.ProductShoes,
		ReceivedAt:  time.Now().UTC(),
	}

	// Act: save
	_, err = repo.Save(ctx, p1)
	require.NoError(t, err)
	_, err = repo.Save(ctx, p2)
	require.NoError(t, err)
	_, err = repo.Save(ctx, p3)
	require.NoError(t, err)

	// Act: list by reception
	list, err := repo.ListByReception(ctx, recID)
	require.NoError(t, err)
	require.Len(t, list, 3)
	require.Equal(t, p1.Type, list[0].Type)
	require.Equal(t, p2.Type, list[1].Type)
	require.Equal(t, p3.Type, list[2].Type)

	// Act: delete last (LIFO)
	deleted, err := repo.DeleteLast(ctx, recID)
	require.NoError(t, err)
	require.NotNil(t, deleted)
	require.Equal(t, p3.Type, deleted.Type)

	// Assert: list after delete
	list, err = repo.ListByReception(ctx, recID)
	require.NoError(t, err)
	require.Len(t, list, 2)
}
