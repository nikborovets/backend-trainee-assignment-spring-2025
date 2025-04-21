package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

// ReceptionRepositoryForClose — интерфейс для работы с приёмками (закрытие)
type ReceptionRepositoryForClose interface {
	GetActive(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error)
	Save(ctx context.Context, reception entities.Reception) (entities.Reception, error)
}

// CloseReceptionUseCase — интерактор для закрытия приёмки
type CloseReceptionUseCase struct {
	repo ReceptionRepositoryForClose
}

func NewCloseReceptionUseCase(repo ReceptionRepositoryForClose) *CloseReceptionUseCase {
	return &CloseReceptionUseCase{repo: repo}
}

// Execute закрывает приёмку, если роль pvz_staff и приёмка открыта
func (uc *CloseReceptionUseCase) Execute(ctx context.Context, user entities.User, pvzID uuid.UUID) (entities.Reception, error) {
	if user.Role != entities.UserRolePVZStaff {
		return entities.Reception{}, errors.New("только сотрудник ПВЗ может закрывать приёмку")
	}
	rec, err := uc.repo.GetActive(ctx, pvzID)
	if err != nil {
		return entities.Reception{}, err
	}
	if rec == nil || !rec.IsOpen() {
		return entities.Reception{}, errors.New("нет открытой приёмки для закрытия")
	}
	if err := rec.Close(); err != nil {
		return entities.Reception{}, err
	}
	return uc.repo.Save(ctx, *rec)
}
