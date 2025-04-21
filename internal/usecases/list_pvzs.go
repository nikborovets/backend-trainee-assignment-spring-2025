package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/interfaces"
)

// PVZRepositoryForList — интерфейс для листинга ПВЗ с фильтрами и пагинацией
type PVZRepositoryForList interface {
	List(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]entities.PVZ, error)
}

// ListPVZsUseCase — интерактор для получения списка ПВЗ с фильтрами и пагинацией
type ListPVZsUseCase struct {
	repo          PVZRepositoryForList
	receptionRepo interfaces.ReceptionRepository
	productRepo   interfaces.ProductRepository
}

func NewListPVZsUseCase(repo PVZRepositoryForList, receptionRepo interfaces.ReceptionRepository, productRepo interfaces.ProductRepository) *ListPVZsUseCase {
	return &ListPVZsUseCase{repo: repo, receptionRepo: receptionRepo, productRepo: productRepo}
}

// Execute возвращает список ПВЗ с фильтрами по дате и пагинацией
func (uc *ListPVZsUseCase) Execute(ctx context.Context, user entities.User, startDate, endDate *time.Time, page, limit int) ([]entities.PVZ, error) {
	if user.Role != entities.UserRolePVZStaff && user.Role != entities.UserRoleModerator {
		return nil, context.Canceled // доступ только для staff/moderator
	}
	return uc.repo.List(ctx, startDate, endDate, page, limit)
}

func (uc *ListPVZsUseCase) GetReceptionsByPVZ(ctx context.Context, pvzID uuid.UUID) ([]entities.Reception, error) {
	return uc.receptionRepo.ListByPVZ(ctx, pvzID)
}

func (uc *ListPVZsUseCase) GetProductsByReception(ctx context.Context, receptionID uuid.UUID) ([]entities.Product, error) {
	return uc.productRepo.ListByReception(ctx, receptionID)
}

// ListPVZsUseCaseIface — интерфейс для моков и контроллеров
type ListPVZsUseCaseIface interface {
	Execute(ctx context.Context, user entities.User, startDate, endDate *time.Time, page, limit int) ([]entities.PVZ, error)
	GetReceptionsByPVZ(ctx context.Context, pvzID uuid.UUID) ([]entities.Reception, error)
	GetProductsByReception(ctx context.Context, receptionID uuid.UUID) ([]entities.Product, error)
}
