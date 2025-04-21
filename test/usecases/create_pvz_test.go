package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
)

type mockPVZRepo struct {
	saveFunc func(ctx context.Context, pvz entities.PVZ) (entities.PVZ, error)
}

func (m *mockPVZRepo) Save(ctx context.Context, pvz entities.PVZ) (entities.PVZ, error) {
	return m.saveFunc(ctx, pvz)
}

func TestCreatePVZUseCase_Execute(t *testing.T) {
	// Arrange
	repo := &mockPVZRepo{
		saveFunc: func(ctx context.Context, pvz entities.PVZ) (entities.PVZ, error) {
			pvz.ID = uuid.New()
			return pvz, nil
		},
	}
	uc := usecases.NewCreatePVZUseCase(repo)
	moderator := entities.User{
		ID:               uuid.New(),
		Role:             entities.UserRoleModerator,
		Email:            "mod@avito.ru",
		RegistrationDate: time.Now(),
	}

	// Act
	pvz, err := uc.Execute(context.Background(), moderator, entities.CityMoscow)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pvz.City != entities.CityMoscow {
		t.Errorf("expected city %s, got %s", entities.CityMoscow, pvz.City)
	}

	// Arrange: невалидный город
	// Act
	_, err = uc.Execute(context.Background(), moderator, "Воронеж")
	// Assert
	if err == nil {
		t.Error("expected error for invalid city")
	}

	// Arrange: не модератор
	client := entities.User{
		ID:               uuid.New(),
		Role:             entities.UserRoleClient,
		Email:            "client@avito.ru",
		RegistrationDate: time.Now(),
	}
	// Act
	_, err = uc.Execute(context.Background(), client, entities.CityMoscow)
	// Assert
	if err == nil {
		t.Error("expected error for non-moderator user")
	}
}
