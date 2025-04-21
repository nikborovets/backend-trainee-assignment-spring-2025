package interfaces

import (
	"context"
)

// ProductController — интерфейс контроллера товаров (см. .puml)
type ProductController interface {
	// AddProduct добавляет товар в приёмку
	AddProduct(ctx context.Context, req AddProductRequest) (ProductDTO, error)
	// DeleteLastProduct удаляет последний товар из приёмки по PVZ
	DeleteLastProduct(ctx context.Context, pvzID string) error
}
