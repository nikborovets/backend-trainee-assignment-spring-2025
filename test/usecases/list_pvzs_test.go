package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
)

type mockPVZRepoForList struct {
	listFunc func(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]entities.PVZ, error)
}

func (m *mockPVZRepoForList) List(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]entities.PVZ, error) {
	return m.listFunc(ctx, startDate, endDate, page, limit)
}

func TestListPVZsUseCase_Execute(t *testing.T) {
	// Arrange
	pvz := entities.PVZ{ID: [16]byte{1}, City: "Москва"}
	repo := &mockPVZRepoForList{
		listFunc: func(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]entities.PVZ, error) {
			return []entities.PVZ{pvz}, nil
		},
	}
	uc := usecases.NewListPVZsUseCase(repo)
	staff := entities.User{Role: entities.UserRolePVZStaff}
	moderator := entities.User{Role: entities.UserRoleModerator}

	// Act
	res, err := uc.Execute(context.Background(), staff, nil, nil, 1, 10)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 || res[0].ID != pvz.ID {
		t.Error("expected one PVZ in result")
	}

	// Act: moderator
	res, err = uc.Execute(context.Background(), moderator, nil, nil, 1, 10)
	if err != nil || len(res) != 1 {
		t.Error("moderator should have access")
	}

	// Act: не staff/moderator
	client := entities.User{Role: entities.UserRoleClient}
	_, err = uc.Execute(context.Background(), client, nil, nil, 1, 10)
	if err == nil {
		t.Error("expected error for non-staff/moderator user")
	}
}
