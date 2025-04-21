package usecases_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockProductRepo struct {
	saveFn func(ctx context.Context, product entities.Product) (entities.Product, error)
}

func (m *mockProductRepo) Save(ctx context.Context, product entities.Product) (entities.Product, error) {
	return m.saveFn(ctx, product)
}

type mockReceptionRepoForAdd struct {
	getActiveFn func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error)
}

func (m *mockReceptionRepoForAdd) GetActive(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
	return m.getActiveFn(ctx, pvzID)
}

func TestAddProductUseCase_Execute(t *testing.T) {
	// Arrange
	pvzID := uuid.New()
	rec := &entities.Reception{ID: pvzID, Status: entities.ReceptionInProgress}
	product := entities.Product{ID: uuid.New(), ReceptionID: pvzID, Type: entities.ProductElectronics}
	user := entities.User{Role: entities.UserRolePVZStaff}

	uc := usecases.NewAddProductUseCase(
		&mockProductRepo{saveFn: func(ctx context.Context, p entities.Product) (entities.Product, error) {
			return product, nil
		}},
		&mockReceptionRepoForAdd{getActiveFn: func(ctx context.Context, id uuid.UUID) (*entities.Reception, error) {
			return rec, nil
		}},
	)

	ctx := context.Background()

	// Act
	res, err := uc.Execute(ctx, user, pvzID, entities.ProductElectronics)

	// Assert
	require.NoError(t, err)
	require.Equal(t, product, res)

	// Не pvz_staff
	user.Role = entities.UserRoleClient
	_, err = uc.Execute(ctx, user, pvzID, entities.ProductElectronics)
	assert.Error(t, err)

	// Некорректный тип
	user.Role = entities.UserRolePVZStaff
	_, err = uc.Execute(ctx, user, pvzID, "еда")
	assert.Error(t, err)

	// Нет открытой приёмки
	uc = usecases.NewAddProductUseCase(
		&mockProductRepo{saveFn: func(ctx context.Context, p entities.Product) (entities.Product, error) {
			return product, nil
		}},
		&mockReceptionRepoForAdd{getActiveFn: func(ctx context.Context, id uuid.UUID) (*entities.Reception, error) {
			return nil, nil
		}},
	)
	_, err = uc.Execute(ctx, user, pvzID, entities.ProductElectronics)
	assert.Error(t, err)
}
