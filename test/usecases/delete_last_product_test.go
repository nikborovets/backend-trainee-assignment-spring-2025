package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockProductRepoForDelete struct {
	deleteFn func(ctx context.Context, productID uuid.UUID) error
}

func (m *mockProductRepoForDelete) Delete(ctx context.Context, productID uuid.UUID) error {
	return m.deleteFn(ctx, productID)
}

type mockReceptionRepoForDelete struct {
	getActiveFn func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error)
	saveFn      func(ctx context.Context, reception entities.Reception) (entities.Reception, error)
}

func (m *mockReceptionRepoForDelete) GetActive(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
	return m.getActiveFn(ctx, pvzID)
}

func (m *mockReceptionRepoForDelete) Save(ctx context.Context, reception entities.Reception) (entities.Reception, error) {
	return m.saveFn(ctx, reception)
}

func TestDeleteLastProductUseCase_Execute(t *testing.T) {
	// Arrange
	pvzID := uuid.New()
	p1, p2 := uuid.New(), uuid.New()
	rec := &entities.Reception{
		ID:       pvzID,
		Status:   entities.ReceptionInProgress,
		Products: []uuid.UUID{p1, p2},
	}
	user := entities.User{
		Role:             entities.UserRolePVZStaff,
		Email:            "staff@avito.ru",
		RegistrationDate: time.Now(),
	}

	productRepo := &mockProductRepoForDelete{
		deleteFn: func(ctx context.Context, productID uuid.UUID) error {
			if productID == p2 {
				return nil
			}
			return assert.AnError
		},
	}
	receptionRepo := &mockReceptionRepoForDelete{
		getActiveFn: func(ctx context.Context, id uuid.UUID) (*entities.Reception, error) {
			return rec, nil
		},
		saveFn: func(ctx context.Context, r entities.Reception) (entities.Reception, error) {
			return *rec, nil
		},
	}
	uc := usecases.NewDeleteLastProductUseCase(productRepo, receptionRepo)
	ctx := context.Background()

	// Act
	err := uc.Execute(ctx, user, pvzID)

	// Assert
	require.NoError(t, err)
	require.Len(t, rec.Products, 1)
	require.Equal(t, p1, rec.Products[0], "LIFO remove failed: wrong product left")

	// Не pvz_staff
	user.Role = entities.UserRoleClient
	err = uc.Execute(ctx, user, pvzID)
	assert.Error(t, err)

	// Нет открытой приёмки
	receptionRepo.getActiveFn = func(ctx context.Context, id uuid.UUID) (*entities.Reception, error) {
		return nil, nil
	}
	user.Role = entities.UserRolePVZStaff
	err = uc.Execute(ctx, user, pvzID)
	assert.Error(t, err)

	// Ошибка удаления
	receptionRepo.getActiveFn = func(ctx context.Context, id uuid.UUID) (*entities.Reception, error) {
		return rec, nil
	}
	productRepo.deleteFn = func(ctx context.Context, productID uuid.UUID) error {
		return assert.AnError
	}
	err = uc.Execute(ctx, user, pvzID)
	assert.Error(t, err)
}
