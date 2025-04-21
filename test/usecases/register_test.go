package usecases_test

import (
	"context"
	"testing"

	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
)

type mockUserRepoForRegister struct {
	createFunc     func(ctx context.Context, user entities.User, passwordHash string) (entities.User, error)
	getByEmailFunc func(ctx context.Context, email string) (*entities.User, error)
}

func (m *mockUserRepoForRegister) Create(ctx context.Context, user entities.User, passwordHash string) (entities.User, error) {
	return m.createFunc(ctx, user, passwordHash)
}
func (m *mockUserRepoForRegister) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	return m.getByEmailFunc(ctx, email)
}

func TestRegisterUseCase_Execute(t *testing.T) {
	// Arrange
	repo := &mockUserRepoForRegister{
		createFunc: func(ctx context.Context, user entities.User, passwordHash string) (entities.User, error) {
			user.ID = [16]byte{1}
			return user, nil
		},
		getByEmailFunc: func(ctx context.Context, email string) (*entities.User, error) {
			return nil, nil
		},
	}
	uc := usecases.NewRegisterUseCase(repo)

	// Act
	user, err := uc.Execute(context.Background(), "test@avito.ru", "pass123", entities.UserRoleClient)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "test@avito.ru" {
		t.Errorf("expected email test@avito.ru, got %s", user.Email)
	}

	// Arrange: дублирующий email
	repo.getByEmailFunc = func(ctx context.Context, email string) (*entities.User, error) {
		u := user
		return &u, nil
	}
	// Act
	_, err = uc.Execute(context.Background(), "test@avito.ru", "pass123", entities.UserRoleClient)
	// Assert
	if err == nil {
		t.Error("expected error for duplicate email")
	}

	// Arrange: невалидная роль
	repo.getByEmailFunc = func(ctx context.Context, email string) (*entities.User, error) { return nil, nil }
	// Act
	_, err = uc.Execute(context.Background(), "test2@avito.ru", "pass123", "admin")
	// Assert
	if err == nil {
		t.Error("expected error for invalid role")
	}

	// Arrange: пустой email/пароль
	_, err = uc.Execute(context.Background(), "", "pass123", entities.UserRoleClient)
	if err == nil {
		t.Error("expected error for empty email")
	}
	_, err = uc.Execute(context.Background(), "test3@avito.ru", "", entities.UserRoleClient)
	if err == nil {
		t.Error("expected error for empty password")
	}
}
