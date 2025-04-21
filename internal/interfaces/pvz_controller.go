package interfaces

import (
	"context"
)

// PVZController — интерфейс контроллера ПВЗ (см. .puml)
type PVZController interface {
	// CreatePVZ создаёт новый ПВЗ
	CreatePVZ(ctx context.Context, req PVZDTO) (PVZDTO, error)
	// ListPVZ возвращает список ПВЗ с фильтрами и пагинацией
	ListPVZ(ctx context.Context, params ListParams) ([]FullPVZDTO, error)
	// CloseLastReception закрывает последнюю приёмку по PVZ
	CloseLastReception(ctx context.Context, pvzID string) (ReceptionDTO, error)
	// DeleteLastProduct удаляет последний товар из приёмки по PVZ
	DeleteLastProduct(ctx context.Context, pvzID string) error
}
