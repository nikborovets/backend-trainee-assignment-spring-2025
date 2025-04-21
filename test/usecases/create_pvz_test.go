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

type mockPVZRepoForCreate struct {
	saveFn func(ctx context.Context, pvz entities.PVZ) (entities.PVZ, error)
}

func (m *mockPVZRepoForCreate) Save(ctx context.Context, pvz entities.PVZ) (entities.PVZ, error) {
	return m.saveFn(ctx, pvz)
}

func TestCreatePVZUseCase_Execute(t *testing.T) {
	// Arrange
	pvz := entities.PVZ{ID: uuid.New(), City: entities.CityMoscow, RegistrationDate: time.Now()}
	user := entities.User{Role: entities.UserRoleModerator}

	repo := &mockPVZRepoForCreate{
		saveFn: func(ctx context.Context, p entities.PVZ) (entities.PVZ, error) {
			return pvz, nil
		},
	}
	uc := usecases.NewCreatePVZUseCase(repo)
	ctx := context.Background()

	// Act
	res, err := uc.Execute(ctx, user, entities.CityMoscow)

	// Assert
	require.NoError(t, err)
	require.Equal(t, entities.CityMoscow, res.City)

	// Некорректный город
	_, err = uc.Execute(ctx, user, "Тверь")
	assert.Error(t, err)

	// Не модератор
	user.Role = entities.UserRoleClient
	_, err = uc.Execute(ctx, user, entities.CityMoscow)
	assert.Error(t, err)
}
