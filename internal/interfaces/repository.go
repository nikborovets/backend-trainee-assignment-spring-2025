package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

// UserRepository — интерфейс для работы с пользователями (см. .puml)
type UserRepository interface {
	Create(ctx context.Context, user entities.User, passwordHash string) (entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, string, error) // user, passwordHash
}

// PVZRepository — интерфейс для работы с ПВЗ
type PVZRepository interface {
	Save(ctx context.Context, pvz entities.PVZ) (entities.PVZ, error)
	List(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]entities.PVZ, error)
}

// ReceptionRepository — интерфейс для работы с приёмками
type ReceptionRepository interface {
	Save(ctx context.Context, reception entities.Reception) (entities.Reception, error)
	CloseLast(ctx context.Context, pvzID uuid.UUID) (entities.Reception, error)
	GetActive(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error)
}

// ProductRepository — интерфейс для работы с товарами
type ProductRepository interface {
	Save(ctx context.Context, product entities.Product) (entities.Product, error)
	DeleteLast(ctx context.Context, pvzID uuid.UUID) error
	ListByReception(ctx context.Context, receptionID uuid.UUID) ([]entities.Product, error)
}
