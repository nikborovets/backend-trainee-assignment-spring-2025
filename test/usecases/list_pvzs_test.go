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

type mockPVZRepoForList struct {
	listFn func(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]entities.PVZ, error)
}

func (m *mockPVZRepoForList) List(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]entities.PVZ, error) {
	return m.listFn(ctx, startDate, endDate, page, limit)
}

func TestListPVZsUseCase_Execute(t *testing.T) {
	// Arrange
	pvz := entities.PVZ{ID: uuid.New(), City: entities.CityMoscow}
	user := entities.User{Role: entities.UserRolePVZStaff}

	repo := &mockPVZRepoForList{
		listFn: func(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]entities.PVZ, error) {
			return []entities.PVZ{pvz}, nil
		},
	}
	uc := usecases.NewListPVZsUseCase(repo)
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
	res, err = uc.Execute(ctx, user, nil, nil, 1, 10)
	assert.Error(t, err)
}
