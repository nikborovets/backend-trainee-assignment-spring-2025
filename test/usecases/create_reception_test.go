package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
)

type mockReceptionRepo struct {
	getActiveFunc func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error)
	saveFunc      func(ctx context.Context, reception entities.Reception) (entities.Reception, error)
}

func (m *mockReceptionRepo) GetActive(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
	return m.getActiveFunc(ctx, pvzID)
}
func (m *mockReceptionRepo) Save(ctx context.Context, reception entities.Reception) (entities.Reception, error) {
	return m.saveFunc(ctx, reception)
}

func TestCreateReceptionUseCase_Execute(t *testing.T) {
	// Arrange
	pvzID := uuid.New()
	repo := &mockReceptionRepo{
		getActiveFunc: func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
			return nil, nil
		},
		saveFunc: func(ctx context.Context, reception entities.Reception) (entities.Reception, error) {
			reception.ID = uuid.New()
			return reception, nil
		},
	}
	uc := usecases.NewCreateReceptionUseCase(repo)
	staff := entities.User{
		ID:               uuid.New(),
		Role:             entities.UserRolePVZStaff,
		Email:            "staff@avito.ru",
		RegistrationDate: time.Now(),
	}

	// Act
	rec, err := uc.Execute(context.Background(), staff, pvzID)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.PVZID != pvzID {
		t.Errorf("expected pvzID %s, got %s", pvzID, rec.PVZID)
	}
	if rec.Status != entities.ReceptionInProgress {
		t.Errorf("expected status in_progress, got %s", rec.Status)
	}

	// Arrange: не pvz_staff
	client := entities.User{
		ID:               uuid.New(),
		Role:             entities.UserRoleClient,
		Email:            "client@avito.ru",
		RegistrationDate: time.Now(),
	}
	// Act
	_, err = uc.Execute(context.Background(), client, pvzID)
	// Assert
	if err == nil {
		t.Error("expected error for non-pvz_staff user")
	}

	// Arrange: уже есть открытая приёмка
	repo.getActiveFunc = func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
		r := &entities.Reception{Status: entities.ReceptionInProgress}
		return r, nil
	}
	// Act
	_, err = uc.Execute(context.Background(), staff, pvzID)
	// Assert
	if err == nil {
		t.Error("expected error if active reception exists")
	}
}
