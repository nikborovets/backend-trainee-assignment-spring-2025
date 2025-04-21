package usecases

import (
	"context"
	"time"

	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

// PVZRepositoryForList — интерфейс для листинга ПВЗ с фильтрами и пагинацией
type PVZRepositoryForList interface {
	List(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]entities.PVZ, error)
}

// ListPVZsUseCase — интерактор для получения списка ПВЗ с фильтрами и пагинацией
type ListPVZsUseCase struct {
	repo PVZRepositoryForList
}

func NewListPVZsUseCase(repo PVZRepositoryForList) *ListPVZsUseCase {
	return &ListPVZsUseCase{repo: repo}
}

// Execute возвращает список ПВЗ с фильтрами по дате и пагинацией
func (uc *ListPVZsUseCase) Execute(ctx context.Context, user entities.User, startDate, endDate *time.Time, page, limit int) ([]entities.PVZ, error) {
	if user.Role != entities.UserRolePVZStaff && user.Role != entities.UserRoleModerator {
		return nil, context.Canceled // доступ только для staff/moderator
	}
	return uc.repo.List(ctx, startDate, endDate, page, limit)
}
