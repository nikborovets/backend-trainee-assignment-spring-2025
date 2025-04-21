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

type mockUserRepoForRegister struct {
	createFn     func(ctx context.Context, user entities.User, passwordHash string) (entities.User, error)
	getByEmailFn func(ctx context.Context, email string) (*entities.User, error)
}

func (m *mockUserRepoForRegister) Create(ctx context.Context, user entities.User, passwordHash string) (entities.User, error) {
	return m.createFn(ctx, user, passwordHash)
}
func (m *mockUserRepoForRegister) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	return m.getByEmailFn(ctx, email)
}

func TestRegisterUseCase_Execute(t *testing.T) {
	// Arrange
	user := entities.User{ID: uuid.New(), Email: "test@avito.ru", Role: entities.UserRoleClient, RegistrationDate: time.Now()}
	repo := &mockUserRepoForRegister{
		createFn: func(ctx context.Context, u entities.User, hash string) (entities.User, error) {
			return user, nil
		},
		getByEmailFn: func(ctx context.Context, email string) (*entities.User, error) {
			return nil, nil
		},
	}
	uc := usecases.NewRegisterUseCase(repo)
	ctx := context.Background()

	// Act
	res, err := uc.Execute(ctx, "test@avito.ru", "password", entities.UserRoleClient)

	// Assert
	require.NoError(t, err)
	require.Equal(t, user.Email, res.Email)

	// Дубликат email
	repo.getByEmailFn = func(ctx context.Context, email string) (*entities.User, error) {
		return &user, nil
	}
	_, err = uc.Execute(ctx, "test@avito.ru", "password", entities.UserRoleClient)
	assert.Error(t, err)

	// Некорректная роль
	repo.getByEmailFn = func(ctx context.Context, email string) (*entities.User, error) {
		return nil, nil
	}
	_, err = uc.Execute(ctx, "test@avito.ru", "password", "hacker")
	assert.Error(t, err)

	// Пустой email
	_, err = uc.Execute(ctx, "", "password", entities.UserRoleClient)
	assert.Error(t, err)

	// Пустой пароль
	_, err = uc.Execute(ctx, "test@avito.ru", "", entities.UserRoleClient)
	assert.Error(t, err)
}
