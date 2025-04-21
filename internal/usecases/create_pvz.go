package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

// PVZRepository — интерфейс для работы с ПВЗ (см. CA_c4_class.puml)
type PVZRepository interface {
	Save(ctx context.Context, pvz entities.PVZ) (entities.PVZ, error)
}

// CreatePVZUseCase — интерактор для создания ПВЗ
// Только модератор может создать ПВЗ, город должен быть разрешён

type CreatePVZUseCase struct {
	pvzRepo PVZRepository
}

func NewCreatePVZUseCase(pvzRepo PVZRepository) *CreatePVZUseCase {
	return &CreatePVZUseCase{pvzRepo: pvzRepo}
}

// Execute создаёт новый ПВЗ, если город разрешён и роль — модератор
func (uc *CreatePVZUseCase) Execute(ctx context.Context, user entities.User, city entities.City) (entities.PVZ, error) {
	if !entities.ValidateUserRole(user.Role) || user.Role != entities.UserRoleModerator {
		return entities.PVZ{}, errors.New("только модератор может создавать ПВЗ")
	}
	if !entities.ValidateCity(city) {
		return entities.PVZ{}, errors.New("ПВЗ можно создать только в Москве, Санкт-Петербурге или Казани")
	}
	pvz := entities.PVZ{
		ID:               entities.GenerateUUID(),
		RegistrationDate: entities.NowUTC(),
		City:             city,
		Receptions:       []uuid.UUID{},
	}
	return uc.pvzRepo.Save(ctx, pvz)
}
