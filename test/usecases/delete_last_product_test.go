package usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
)

type mockProductRepoForDelete struct {
	deleteFunc func(ctx context.Context, productID uuid.UUID) error
}

func (m *mockProductRepoForDelete) Delete(ctx context.Context, productID uuid.UUID) error {
	return m.deleteFunc(ctx, productID)
}

type mockReceptionRepoForDelete struct {
	getActiveFunc func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error)
	saveFunc      func(ctx context.Context, reception entities.Reception) (entities.Reception, error)
}

func (m *mockReceptionRepoForDelete) GetActive(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
	return m.getActiveFunc(ctx, pvzID)
}
func (m *mockReceptionRepoForDelete) Save(ctx context.Context, reception entities.Reception) (entities.Reception, error) {
	return m.saveFunc(ctx, reception)
}

func TestDeleteLastProductUseCase_Execute(t *testing.T) {
	// Arrange
	pvzID := uuid.New()
	p1, p2 := uuid.New(), uuid.New()
	rec := &entities.Reception{
		ID:       uuid.New(),
		PVZID:    pvzID,
		Status:   entities.ReceptionInProgress,
		Products: []uuid.UUID{p1, p2},
	}
	productRepo := &mockProductRepoForDelete{
		deleteFunc: func(ctx context.Context, productID uuid.UUID) error {
			if productID == p2 {
				return nil
			}
			return errors.New("wrong product deleted")
		},
	}
	receptionRepo := &mockReceptionRepoForDelete{
		getActiveFunc: func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
			return rec, nil
		},
		saveFunc: func(ctx context.Context, reception entities.Reception) (entities.Reception, error) {
			return reception, nil
		},
	}
	uc := usecases.NewDeleteLastProductUseCase(productRepo, receptionRepo)
	staff := entities.User{
		ID:               uuid.New(),
		Role:             entities.UserRolePVZStaff,
		Email:            "staff@avito.ru",
		RegistrationDate: time.Now(),
	}

	// Act
	err := uc.Execute(context.Background(), staff, pvzID)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.Products) != 1 || rec.Products[0] != p1 {
		t.Error("LIFO remove failed: wrong product left")
	}

	// Arrange: не pvz_staff
	client := entities.User{
		ID:               uuid.New(),
		Role:             entities.UserRoleClient,
		Email:            "client@avito.ru",
		RegistrationDate: time.Now(),
	}
	// Act
	err = uc.Execute(context.Background(), client, pvzID)
	// Assert
	if err == nil {
		t.Error("expected error for non-pvz_staff user")
	}

	// Arrange: нет открытой приёмки
	receptionRepo.getActiveFunc = func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
		return nil, nil
	}
	// Act
	err = uc.Execute(context.Background(), staff, pvzID)
	// Assert
	if err == nil {
		t.Error("expected error if no open reception")
	}

	// Arrange: ошибка удаления товара
	rec.Products = []uuid.UUID{p1}
	productRepo.deleteFunc = func(ctx context.Context, productID uuid.UUID) error {
		return errors.New("delete failed")
	}
	receptionRepo.getActiveFunc = func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
		return rec, nil
	}
	// Act
	err = uc.Execute(context.Background(), staff, pvzID)
	// Assert
	if err == nil || err.Error() != "delete failed" {
		t.Error("expected error from productRepo.Delete")
	}
}
