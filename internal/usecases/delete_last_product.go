package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

// ProductRepositoryForDelete — интерфейс для удаления товара (LIFO)
type ProductRepositoryForDelete interface {
	DeleteLast(ctx context.Context, receptionID uuid.UUID) (*entities.Product, error)
}

type ReceptionRepositoryForDelete interface {
	GetActive(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error)
}

// DeleteLastProductUseCase — интерактор для удаления последнего товара из приёмки (LIFO)
type DeleteLastProductUseCase struct {
	productRepo   ProductRepositoryForDelete
	receptionRepo ReceptionRepositoryForDelete
}

func NewDeleteLastProductUseCase(productRepo ProductRepositoryForDelete, receptionRepo ReceptionRepositoryForDelete) *DeleteLastProductUseCase {
	return &DeleteLastProductUseCase{productRepo: productRepo, receptionRepo: receptionRepo}
}

// Execute удаляет последний товар из незакрытой приёмки, если роль pvz_staff
func (uc *DeleteLastProductUseCase) Execute(ctx context.Context, user entities.User, pvzID uuid.UUID) error {
	if user.Role != entities.UserRolePVZStaff {
		return errors.New("только сотрудник ПВЗ может удалять товары")
	}

	rec, err := uc.receptionRepo.GetActive(ctx, pvzID)
	if err != nil {
		return err
	}
	if rec == nil || !rec.IsOpen() {
		return errors.New("нет открытой приёмки для удаления товара")
	}

	// Удаляем последний товар через репозиторий
	product, err := uc.productRepo.DeleteLast(ctx, rec.ID)
	if err != nil {
		return err
	}

	if product == nil {
		return errors.New("нет товаров для удаления")
	}

	return nil
}

// DeleteLastProductUseCaseIface — интерфейс для моков и контроллеров
type DeleteLastProductUseCaseIface interface {
	Execute(ctx context.Context, user entities.User, pvzID uuid.UUID) error
}
