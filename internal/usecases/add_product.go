package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

// ProductRepository — интерфейс для работы с товарами (см. CA_c4_class.puml)
type ProductRepository interface {
	Save(ctx context.Context, product entities.Product) (entities.Product, error)
}

// ReceptionRepositoryForAdd — интерфейс для получения незакрытой приёмки
// (можно использовать тот же, что и для CreateReception, но для явности)
type ReceptionRepositoryForAdd interface {
	GetActive(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error)
}

// AddProductUseCase — интерактор для добавления товара в приёмку
// Только pvz_staff, только в незакрытую приёмку, тип товара валиден

type AddProductUseCase struct {
	productRepo   ProductRepository
	receptionRepo ReceptionRepositoryForAdd
}

func NewAddProductUseCase(productRepo ProductRepository, receptionRepo ReceptionRepositoryForAdd) *AddProductUseCase {
	return &AddProductUseCase{productRepo: productRepo, receptionRepo: receptionRepo}
}

// Execute добавляет товар в незакрытую приёмку, если роль pvz_staff и тип валиден
func (uc *AddProductUseCase) Execute(ctx context.Context, user entities.User, pvzID uuid.UUID, productType entities.ProductType) (entities.Product, error) {
	if user.Role != entities.UserRolePVZStaff {
		return entities.Product{}, errors.New("только сотрудник ПВЗ может добавлять товары")
	}
	if !entities.ValidateProductType(productType) {
		return entities.Product{}, errors.New("некорректный тип товара")
	}
	rec, err := uc.receptionRepo.GetActive(ctx, pvzID)
	if err != nil {
		return entities.Product{}, err
	}
	if rec == nil || !rec.IsOpen() {
		return entities.Product{}, errors.New("нет открытой приёмки для добавления товара")
	}
	product := entities.Product{
		ID:          uuid.New(),
		ReceptionID: rec.ID,
		Type:        productType,
		ReceivedAt:  time.Now().UTC(),
	}
	return uc.productRepo.Save(ctx, product)
}
