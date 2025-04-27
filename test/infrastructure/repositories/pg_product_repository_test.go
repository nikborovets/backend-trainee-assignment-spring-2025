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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// Загружаем переменные окружения из .env файла
	_ = godotenv.Load("../../../.env")

	// Вывод значения для отладки
	dsn := os.Getenv("TEST_PG_DSN")
	if dsn != "" {
		println("TEST_PG_DSN loaded:", dsn)
	}
}

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
		DateTime:    time.Now().Add(-2 * time.Hour).UTC(),
	}
	p2 := entities.Product{
		ID:          uuid.New(),
		ReceptionID: recID,
		Type:        entities.ProductClothes,
		DateTime:    time.Now().Add(-1 * time.Hour).UTC(),
	}
	p3 := entities.Product{
		ID:          uuid.New(),
		ReceptionID: recID,
		Type:        entities.ProductShoes,
		DateTime:    time.Now().UTC(),
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

// TestPGProductRepository_ListByReception проверяет получение всех товаров по приёмке
func TestPGProductRepository_ListByReception(t *testing.T) {
	db := setupProductTestDB(t)
	repo := repositories.NewPGProductRepository(db)

	ctx := context.Background()

	// Создаем ПВЗ перед созданием приёмки (из-за внешнего ключа)
	pvzID := uuid.New()
	_, err := db.Exec(`INSERT INTO pvz (id, registration_date, city) VALUES ($1, $2, $3)`,
		pvzID, time.Now().UTC(), "Москва")
	require.NoError(t, err)

	// Создаем приёмку перед созданием товаров (из-за внешнего ключа)
	recID := uuid.New()
	_, err = db.Exec(`INSERT INTO reception (id, pvz_id, status, date_time) VALUES ($1, $2, $3, $4)`,
		recID, pvzID, "in_progress", time.Now().UTC())
	require.NoError(t, err)

	// Добавляем несколько товаров
	products := []entities.Product{
		{
			ID:          uuid.New(),
			ReceptionID: recID,
			Type:        entities.ProductElectronics,
			DateTime:    time.Now().Add(-2 * time.Hour).UTC(),
		},
		{
			ID:          uuid.New(),
			ReceptionID: recID,
			Type:        entities.ProductClothes,
			DateTime:    time.Now().Add(-1 * time.Hour).UTC(),
		},
		{
			ID:          uuid.New(),
			ReceptionID: recID,
			Type:        entities.ProductShoes,
			DateTime:    time.Now().UTC(),
		},
	}

	// Сохраняем в БД
	for _, p := range products {
		_, err := repo.Save(ctx, p)
		require.NoError(t, err)
	}

	// Получаем все товары для приёмки
	found, err := repo.ListByReception(ctx, recID)
	require.NoError(t, err)
	require.Len(t, found, 3)

	// Проверяем сортировку по времени (должны быть в том же порядке)
	assert.Equal(t, products[0].ID, found[0].ID)
	assert.Equal(t, products[1].ID, found[1].ID)
	assert.Equal(t, products[2].ID, found[2].ID)
}
