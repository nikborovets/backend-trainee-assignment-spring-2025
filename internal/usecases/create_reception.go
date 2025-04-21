package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

// ReceptionRepository — интерфейс для работы с приёмками (см. CA_c4_class.puml)
type ReceptionRepository interface {
	GetActive(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error)
	Save(ctx context.Context, reception entities.Reception) (entities.Reception, error)
}

// CreateReceptionUseCase — интерактор для создания приёмки
// Только pvz_staff может создать приёмку, на PVZ может быть только одна открытая приёмка

type CreateReceptionUseCase struct {
	repo ReceptionRepository
}

func NewCreateReceptionUseCase(repo ReceptionRepository) *CreateReceptionUseCase {
	return &CreateReceptionUseCase{repo: repo}
}

// Execute создаёт новую приёмку, если нет открытой приёмки на PVZ и роль — pvz_staff
func (uc *CreateReceptionUseCase) Execute(ctx context.Context, user entities.User, pvzID uuid.UUID) (entities.Reception, error) {
	if user.Role != "pvz_staff" {
		return entities.Reception{}, errors.New("только сотрудник ПВЗ может создавать приёмку")
	}
	active, err := uc.repo.GetActive(ctx, pvzID)
	if err != nil {
		return entities.Reception{}, err
	}
	if active != nil && active.IsOpen() {
		return entities.Reception{}, errors.New("у ПВЗ уже есть открытая приёмка")
	}
	rec := entities.Reception{
		ID:       uuid.New(),
		PVZID:    pvzID,
		Products: []uuid.UUID{},
		Status:   entities.ReceptionInProgress,
		DateTime: time.Now().UTC(),
	}
	return uc.repo.Save(ctx, rec)
}
