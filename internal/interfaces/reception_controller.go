package interfaces

import (
	"context"
)

// ReceptionController — интерфейс контроллера приёмок (см. .puml)
type ReceptionController interface {
	// CreateReception создаёт новую приёмку
	CreateReception(ctx context.Context, pvzID string) (ReceptionDTO, error)
	// CloseLastReception закрывает последнюю приёмку по PVZ
	CloseLastReception(ctx context.Context, pvzID string) (ReceptionDTO, error)
}
