package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
)

type mockProductRepo struct {
	saveFunc func(ctx context.Context, product entities.Product) (entities.Product, error)
}

func (m *mockProductRepo) Save(ctx context.Context, product entities.Product) (entities.Product, error) {
	return m.saveFunc(ctx, product)
}

type mockReceptionRepoForAdd struct {
	getActiveFunc func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error)
}

func (m *mockReceptionRepoForAdd) GetActive(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
	return m.getActiveFunc(ctx, pvzID)
}

func TestAddProductUseCase_Execute(t *testing.T) {
	// Arrange
	pvzID := uuid.New()
	rec := &entities.Reception{
		ID:     uuid.New(),
		PVZID:  pvzID,
		Status: entities.ReceptionInProgress,
	}
	productRepo := &mockProductRepo{
		saveFunc: func(ctx context.Context, product entities.Product) (entities.Product, error) {
			product.ID = uuid.New()
			return product, nil
		},
	}
	receptionRepo := &mockReceptionRepoForAdd{
		getActiveFunc: func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
			return rec, nil
		},
	}
	uc := usecases.NewAddProductUseCase(productRepo, receptionRepo)
	staff := entities.User{
		ID:               uuid.New(),
		Role:             entities.UserRolePVZStaff,
		Email:            "staff@avito.ru",
		RegistrationDate: time.Now(),
	}

	// Act
	product, err := uc.Execute(context.Background(), staff, pvzID, entities.ProductElectronics)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if product.Type != entities.ProductElectronics {
		t.Errorf("expected type %s, got %s", entities.ProductElectronics, product.Type)
	}
	if product.ReceptionID != rec.ID {
		t.Errorf("expected receptionID %s, got %s", rec.ID, product.ReceptionID)
	}

	// Arrange: не pvz_staff
	client := entities.User{
		ID:               uuid.New(),
		Role:             entities.UserRoleClient,
		Email:            "client@avito.ru",
		RegistrationDate: time.Now(),
	}
	// Act
	_, err = uc.Execute(context.Background(), client, pvzID, entities.ProductElectronics)
	// Assert
	if err == nil {
		t.Error("expected error for non-pvz_staff user")
	}

	// Arrange: невалидный тип
	// Act
	_, err = uc.Execute(context.Background(), staff, pvzID, "еда")
	// Assert
	if err == nil {
		t.Error("expected error for invalid product type")
	}

	// Arrange: нет открытой приёмки
	receptionRepo.getActiveFunc = func(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
		return nil, nil
	}
	// Act
	_, err = uc.Execute(context.Background(), staff, pvzID, entities.ProductElectronics)
	// Assert
	if err == nil {
		t.Error("expected error if no open reception")
	}
}
