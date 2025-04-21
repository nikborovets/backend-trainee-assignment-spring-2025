package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

// ProductRepositoryForDelete — интерфейс для удаления товара (LIFO)
type ProductRepositoryForDelete interface {
	Delete(ctx context.Context, productID uuid.UUID) error
}

type ReceptionRepositoryForDelete interface {
	GetActive(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error)
	Save(ctx context.Context, reception entities.Reception) (entities.Reception, error)
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
	lastID, err := rec.RemoveLastProduct()
	if err != nil {
		return err
	}
	if err := uc.productRepo.Delete(ctx, lastID); err != nil {
		return err
	}
	_, err = uc.receptionRepo.Save(ctx, *rec)
	return err
}

// DeleteLastProductUseCaseIface — интерфейс для моков и контроллеров
type DeleteLastProductUseCaseIface interface {
	Execute(ctx context.Context, user entities.User, pvzID uuid.UUID) error
}
