package usecases_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockReceptionRepo struct{ mock.Mock }

func (m *mockReceptionRepo) ListByPVZ(ctx context.Context, pvzID uuid.UUID) ([]entities.Reception, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).([]entities.Reception), args.Error(1)
}

// остальные методы не нужны для этого теста
func (m *mockReceptionRepo) Save(context.Context, entities.Reception) (entities.Reception, error) {
	return entities.Reception{}, nil
}
func (m *mockReceptionRepo) CloseLast(context.Context, uuid.UUID) (entities.Reception, error) {
	return entities.Reception{}, nil
}
func (m *mockReceptionRepo) GetActive(context.Context, uuid.UUID) (*entities.Reception, error) {
	return nil, nil
}

type mockProductRepo struct{ mock.Mock }

func (m *mockProductRepo) ListByReception(ctx context.Context, receptionID uuid.UUID) ([]entities.Product, error) {
	args := m.Called(ctx, receptionID)
	return args.Get(0).([]entities.Product), args.Error(1)
}

// остальные методы не нужны для этого теста
func (m *mockProductRepo) Save(context.Context, entities.Product) (entities.Product, error) {
	return entities.Product{}, nil
}
func (m *mockProductRepo) DeleteLast(context.Context, uuid.UUID) error { return nil }

func TestListPVZsUseCase_GetReceptionsByPVZ(t *testing.T) {
	ctx := context.Background()
	pvzID := uuid.New()
	recs := []entities.Reception{{ID: uuid.New()}, {ID: uuid.New()}}
	repo := new(mockReceptionRepo)
	repo.On("ListByPVZ", ctx, pvzID).Return(recs, nil)
	uc := usecases.NewListPVZsUseCase(nil, repo, nil)

	// Act
	got, err := uc.GetReceptionsByPVZ(ctx, pvzID)

	// Assert
	require.NoError(t, err)
	require.Equal(t, recs, got)
	repo.AssertExpectations(t)
}

func TestListPVZsUseCase_GetProductsByReception(t *testing.T) {
	ctx := context.Background()
	recID := uuid.New()
	products := []entities.Product{{ID: uuid.New()}, {ID: uuid.New()}}
	repo := new(mockProductRepo)
	repo.On("ListByReception", ctx, recID).Return(products, nil)
	uc := usecases.NewListPVZsUseCase(nil, nil, repo)

	// Act
	got, err := uc.GetProductsByReception(ctx, recID)

	// Assert
	require.NoError(t, err)
	require.Equal(t, products, got)
	repo.AssertExpectations(t)
}
