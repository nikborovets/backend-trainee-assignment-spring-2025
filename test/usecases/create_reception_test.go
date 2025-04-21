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

type mockReceptionRepoForCreate struct {
	getActiveFn func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error)
	saveFn      func(ctx context.Context, reception entities.Reception) (entities.Reception, error)
}

func (m *mockReceptionRepoForCreate) GetActive(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
	return m.getActiveFn(ctx, pvzID)
}
func (m *mockReceptionRepoForCreate) Save(ctx context.Context, reception entities.Reception) (entities.Reception, error) {
	return m.saveFn(ctx, reception)
}

func TestCreateReceptionUseCase_Execute(t *testing.T) {
	// Arrange
	pvzID := uuid.New()
	rec := entities.Reception{ID: uuid.New(), PVZID: pvzID, Status: entities.ReceptionInProgress}
	user := entities.User{Role: entities.UserRolePVZStaff}

	repo := &mockReceptionRepoForCreate{
		getActiveFn: func(ctx context.Context, id uuid.UUID) (*entities.Reception, error) {
			return nil, nil
		},
		saveFn: func(ctx context.Context, r entities.Reception) (entities.Reception, error) {
			return rec, nil
		},
	}
	uc := usecases.NewCreateReceptionUseCase(repo)
	ctx := context.Background()

	// Act
	result, err := uc.Execute(ctx, user, pvzID)

	// Assert
	require.NoError(t, err)
	require.Equal(t, rec.PVZID, result.PVZID)
	require.Equal(t, entities.ReceptionInProgress, result.Status)

	// Не pvz_staff
	user.Role = entities.UserRoleClient
	_, err = uc.Execute(ctx, user, pvzID)
	assert.Error(t, err)

	// Уже есть открытая приёмка
	repo.getActiveFn = func(ctx context.Context, id uuid.UUID) (*entities.Reception, error) {
		r := &entities.Reception{Status: entities.ReceptionInProgress}
		return r, nil
	}
	user.Role = entities.UserRolePVZStaff
	_, err = uc.Execute(ctx, user, pvzID)
	assert.Error(t, err)
}
