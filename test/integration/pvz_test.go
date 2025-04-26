package integration_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/configs"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/infrastructure/repositories"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
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

func TestPVZIntegration(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Инициализация репозиториев
	pvzRepo := repositories.NewPGPVZRepository(db)
	receptionRepo := repositories.NewPGReceptionRepository(db)
	productRepo := repositories.NewPGProductRepository(db)

	// Инициализация use cases
	createPVZUC := usecases.NewCreatePVZUseCase(pvzRepo)
	createReceptionUC := usecases.NewCreateReceptionUseCase(receptionRepo)
	addProductUC := usecases.NewAddProductUseCase(productRepo, receptionRepo)
	closeReceptionUC := usecases.NewCloseReceptionUseCase(receptionRepo)

	// 1. Создание ПВЗ
	moderator := entities.User{
		ID:    uuid.New(),
		Email: "moderator@avito.ru",
		Role:  entities.UserRoleModerator,
	}
	pvz, err := createPVZUC.Execute(ctx, moderator, entities.CityMoscow)
	require.NoError(t, err)
	require.NotNil(t, pvz)

	// 2. Создание приёмки
	staff := entities.User{
		ID:    uuid.New(),
		Email: "staff@avito.ru",
		Role:  entities.UserRolePVZStaff,
	}
	reception, err := createReceptionUC.Execute(ctx, staff, pvz.ID)
	require.NoError(t, err)
	require.NotNil(t, reception)
	require.Equal(t, entities.ReceptionInProgress, reception.Status)

	// 3. Добавление 50 товаров
	for i := 0; i < 50; i++ {
		productType := entities.ProductElectronics
		if i%3 == 0 {
			productType = entities.ProductClothes
		} else if i%3 == 1 {
			productType = entities.ProductShoes
		}
		product, err := addProductUC.Execute(ctx, staff, pvz.ID, productType)
		require.NoError(t, err)
		require.NotNil(t, product)
		require.Equal(t, productType, product.Type)
		require.Equal(t, reception.ID, product.ReceptionID)
	}

	// 4. Закрытие приёмки
	closedReception, err := closeReceptionUC.Execute(ctx, staff, pvz.ID)
	require.NoError(t, err)
	require.NotNil(t, closedReception)
	require.Equal(t, entities.ReceptionClosed, closedReception.Status)

	// Проверка, что все товары сохранились
	products, err := productRepo.ListByReception(ctx, reception.ID)
	require.NoError(t, err)
	require.Len(t, products, 50)
}
