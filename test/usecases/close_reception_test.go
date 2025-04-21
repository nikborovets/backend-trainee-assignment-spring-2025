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

type mockReceptionRepoForClose struct {
	getActiveFn func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error)
	saveFn      func(ctx context.Context, reception entities.Reception) (entities.Reception, error)
}

func (m *mockReceptionRepoForClose) GetActive(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
	return m.getActiveFn(ctx, pvzID)
}
func (m *mockReceptionRepoForClose) Save(ctx context.Context, reception entities.Reception) (entities.Reception, error) {
	return m.saveFn(ctx, reception)
}

func TestCloseReceptionUseCase_Execute(t *testing.T) {
	// Arrange
	pvzID := uuid.New()
	rec := entities.Reception{ID: uuid.New(), PVZID: pvzID, Status: entities.ReceptionInProgress}
	user := entities.User{Role: entities.UserRolePVZStaff}

	repo := &mockReceptionRepoForClose{
		getActiveFn: func(ctx context.Context, id uuid.UUID) (*entities.Reception, error) {
			return &rec, nil
		},
		saveFn: func(ctx context.Context, r entities.Reception) (entities.Reception, error) {
			return r, nil
		},
	}
	uc := usecases.NewCloseReceptionUseCase(repo)
	ctx := context.Background()

	// Act
	closed, err := uc.Execute(ctx, user, pvzID)

	// Assert
	require.NoError(t, err)
	require.Equal(t, entities.ReceptionClosed, closed.Status)

	// Не pvz_staff
	user.Role = entities.UserRoleClient
	_, err = uc.Execute(ctx, user, pvzID)
	assert.Error(t, err)

	// Нет открытой приёмки
	repo.getActiveFn = func(ctx context.Context, id uuid.UUID) (*entities.Reception, error) {
		return nil, nil
	}
	user.Role = entities.UserRolePVZStaff
	_, err = uc.Execute(ctx, user, pvzID)
	assert.Error(t, err)

	// Уже закрыта
	repo.getActiveFn = func(ctx context.Context, id uuid.UUID) (*entities.Reception, error) {
		r := &entities.Reception{Status: entities.ReceptionClosed}
		return r, nil
	}
	_, err = uc.Execute(ctx, user, pvzID)
	assert.Error(t, err)
}
