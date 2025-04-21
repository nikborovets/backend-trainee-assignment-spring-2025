package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
)

type mockReceptionRepoForClose struct {
	getActiveFunc func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error)
	saveFunc      func(ctx context.Context, reception entities.Reception) (entities.Reception, error)
}

func (m *mockReceptionRepoForClose) GetActive(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
	return m.getActiveFunc(ctx, pvzID)
}
func (m *mockReceptionRepoForClose) Save(ctx context.Context, reception entities.Reception) (entities.Reception, error) {
	return m.saveFunc(ctx, reception)
}

func TestCloseReceptionUseCase_Execute(t *testing.T) {
	// Arrange
	pvzID := uuid.New()
	rec := &entities.Reception{
		ID:     uuid.New(),
		PVZID:  pvzID,
		Status: entities.ReceptionInProgress,
	}
	repo := &mockReceptionRepoForClose{
		getActiveFunc: func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
			return rec, nil
		},
		saveFunc: func(ctx context.Context, reception entities.Reception) (entities.Reception, error) {
			return reception, nil
		},
	}
	uc := usecases.NewCloseReceptionUseCase(repo)
	staff := entities.User{
		ID:               uuid.New(),
		Role:             entities.UserRolePVZStaff,
		Email:            "staff@avito.ru",
		RegistrationDate: time.Now(),
	}

	// Act
	closed, err := uc.Execute(context.Background(), staff, pvzID)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if closed.Status != entities.ReceptionClosed {
		t.Errorf("expected status closed, got %s", closed.Status)
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

	// Arrange: нет открытой приёмки
	repo.getActiveFunc = func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
		return nil, nil
	}
	// Act
	_, err = uc.Execute(context.Background(), staff, pvzID)
	// Assert
	if err == nil {
		t.Error("expected error if no open reception")
	}

	// Arrange: приёмка уже закрыта
	rec.Status = entities.ReceptionClosed
	repo.getActiveFunc = func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
		return rec, nil
	}
	// Act
	_, err = uc.Execute(context.Background(), staff, pvzID)
	// Assert
	if err == nil {
		t.Error("expected error if reception already closed")
	}
}
