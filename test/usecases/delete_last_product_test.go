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
	deleteLastFn func(ctx context.Context, receptionID uuid.UUID) (*entities.Product, error)
}

func (m *mockProductRepoForDelete) DeleteLast(ctx context.Context, receptionID uuid.UUID) (*entities.Product, error) {
	return m.deleteLastFn(ctx, receptionID)
}

type mockReceptionRepoForDelete struct {
	getActiveFn func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error)
}

func (m *mockReceptionRepoForDelete) GetActive(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
	return m.getActiveFn(ctx, pvzID)
}

func TestDeleteLastProductUseCase_Execute(t *testing.T) {
	// Arrange
	pvzID := uuid.New()
	receptionID := uuid.New()

	activeReception := &entities.Reception{
		ID:       receptionID,
		PVZID:    pvzID,
		Status:   entities.ReceptionInProgress,
		DateTime: time.Now(),
	}

	lastProduct := &entities.Product{
		ID:          uuid.New(),
		ReceptionID: receptionID,
		Type:        entities.ProductElectronics,
		DateTime:    time.Now(),
	}

	user := entities.User{
		Role:             entities.UserRolePVZStaff,
		Email:            "staff@avito.ru",
		RegistrationDate: time.Now(),
	}

	productRepo := &mockProductRepoForDelete{
		deleteLastFn: func(ctx context.Context, recID uuid.UUID) (*entities.Product, error) {
			if recID == receptionID {
				return lastProduct, nil
			}
			return nil, assert.AnError
		},
	}

	receptionRepo := &mockReceptionRepoForDelete{
		getActiveFn: func(ctx context.Context, id uuid.UUID) (*entities.Reception, error) {
			if id == pvzID {
				return activeReception, nil
			}
			return nil, nil
		},
	}

	uc := usecases.NewDeleteLastProductUseCase(productRepo, receptionRepo)
	ctx := context.Background()

	// Act
	err := uc.Execute(ctx, user, pvzID)

	// Assert
	require.NoError(t, err)

	// Проверка, что метод GetActive вызывается корректно
	receptionRepo.getActiveFn = func(ctx context.Context, id uuid.UUID) (*entities.Reception, error) {
		assert.Equal(t, pvzID, id)
		return activeReception, nil
	}

	// Проверка, что метод DeleteLast вызывается с правильным reception_id
	productRepo.deleteLastFn = func(ctx context.Context, recID uuid.UUID) (*entities.Product, error) {
		assert.Equal(t, receptionID, recID)
		return lastProduct, nil
	}

	err = uc.Execute(ctx, user, pvzID)
	require.NoError(t, err)

	// Тест: Не pvz_staff
	user.Role = entities.UserRoleClient
	err = uc.Execute(ctx, user, pvzID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "только сотрудник ПВЗ может удалять товары")

	// Тест: Нет открытой приёмки
	user.Role = entities.UserRolePVZStaff
	receptionRepo.getActiveFn = func(ctx context.Context, id uuid.UUID) (*entities.Reception, error) {
		return nil, nil
	}
	err = uc.Execute(ctx, user, pvzID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "нет открытой приёмки")

	// Тест: Закрытая приемка
	closedReception := &entities.Reception{
		ID:       receptionID,
		PVZID:    pvzID,
		Status:   entities.ReceptionClosed,
		DateTime: time.Now(),
	}
	receptionRepo.getActiveFn = func(ctx context.Context, id uuid.UUID) (*entities.Reception, error) {
		return closedReception, nil
	}
	err = uc.Execute(ctx, user, pvzID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "нет открытой приёмки")

	// Тест: Ошибка при удалении последнего товара
	receptionRepo.getActiveFn = func(ctx context.Context, id uuid.UUID) (*entities.Reception, error) {
		return activeReception, nil
	}
	productRepo.deleteLastFn = func(ctx context.Context, recID uuid.UUID) (*entities.Product, error) {
		return nil, assert.AnError
	}
	err = uc.Execute(ctx, user, pvzID)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)

	// Тест: Нет товаров для удаления
	productRepo.deleteLastFn = func(ctx context.Context, recID uuid.UUID) (*entities.Product, error) {
		return nil, nil
	}
	err = uc.Execute(ctx, user, pvzID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "нет товаров для удаления")
}
