package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Моки для тестирования Execute
type mockPVZRepoForList struct {
	listFn func(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]entities.PVZ, error)
}

func (m *mockPVZRepoForList) List(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]entities.PVZ, error) {
	return m.listFn(ctx, startDate, endDate, page, limit)
}

// Моки для тестирования GetReceptionsByPVZ и GetProductsByReception
type mockReceptionRepoForList struct{ mock.Mock }

func (m *mockReceptionRepoForList) ListByPVZ(ctx context.Context, pvzID uuid.UUID) ([]entities.Reception, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).([]entities.Reception), args.Error(1)
}

// остальные методы не нужны для этого теста
func (m *mockReceptionRepoForList) Save(context.Context, entities.Reception) (entities.Reception, error) {
	return entities.Reception{}, nil
}
func (m *mockReceptionRepoForList) CloseLast(context.Context, uuid.UUID) (entities.Reception, error) {
	return entities.Reception{}, nil
}
func (m *mockReceptionRepoForList) GetActive(context.Context, uuid.UUID) (*entities.Reception, error) {
	return nil, nil
}

type mockProductRepoForList struct{ mock.Mock }

func (m *mockProductRepoForList) ListByReception(ctx context.Context, receptionID uuid.UUID) ([]entities.Product, error) {
	args := m.Called(ctx, receptionID)
	return args.Get(0).([]entities.Product), args.Error(1)
}

// остальные методы не нужны для этого теста
func (m *mockProductRepoForList) Save(context.Context, entities.Product) (entities.Product, error) {
	return entities.Product{}, nil
}
func (m *mockProductRepoForList) DeleteLast(context.Context, uuid.UUID) error { return nil }

func TestListPVZsUseCase_Execute(t *testing.T) {
	// Arrange
	pvz := entities.PVZ{ID: uuid.New(), City: entities.CityMoscow}
	user := entities.User{Role: entities.UserRolePVZStaff}

	repo := &mockPVZRepoForList{
		listFn: func(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]entities.PVZ, error) {
			return []entities.PVZ{pvz}, nil
		},
	}
	uc := usecases.NewListPVZsUseCase(repo, nil, nil)
	ctx := context.Background()

	// Act
	res, err := uc.Execute(ctx, user, nil, nil, 1, 10)

	// Assert
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Equal(t, pvz.ID, res[0].ID)

	// Модератор тоже может
	user.Role = entities.UserRoleModerator
	res, err = uc.Execute(ctx, user, nil, nil, 1, 10)
	require.NoError(t, err)
	require.Len(t, res, 1)

	// Не staff/moderator
	user.Role = "hacker"
	_, err = uc.Execute(ctx, user, nil, nil, 1, 10)
	assert.Error(t, err)
}

func TestListPVZsUseCase_GetReceptionsByPVZ(t *testing.T) {
	ctx := context.Background()
	pvzID := uuid.New()
	recs := []entities.Reception{{ID: uuid.New()}, {ID: uuid.New()}}
	repo := new(mockReceptionRepoForList)
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
	repo := new(mockProductRepoForList)
	repo.On("ListByReception", ctx, recID).Return(products, nil)
	uc := usecases.NewListPVZsUseCase(nil, nil, repo)

	// Act
	got, err := uc.GetProductsByReception(ctx, recID)

	// Assert
	require.NoError(t, err)
	require.Equal(t, products, got)
	repo.AssertExpectations(t)
}
