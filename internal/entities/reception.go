package entities

import (
	"time"

	"errors"

	"github.com/google/uuid"
)

// Reception — приёмка товаров в ПВЗ
// id — UUID
// pvzId — UUID
// products — список товаров (UUID)
// status — in_progress/close
// dateTime — дата и время приёмки

type ReceptionStatus string

const (
	ReceptionInProgress ReceptionStatus = "in_progress"
	ReceptionClosed     ReceptionStatus = "close"
)

type Reception struct {
	ID       uuid.UUID       `json:"id"`
	PVZID    uuid.UUID       `json:"pvzId"`
	Products []uuid.UUID     `json:"products"`
	Status   ReceptionStatus `json:"status"`
	DateTime time.Time       `json:"dateTime"`
}

// Проверяет, открыта ли приёмка
func (r *Reception) IsOpen() bool {
	return r.Status == ReceptionInProgress
}

// Добавляет товар в приёмку, если она открыта
func (r *Reception) AddProduct(productID uuid.UUID) error {
	if !r.IsOpen() {
		return errors.New("приёмка закрыта, нельзя добавить товар")
	}
	r.Products = append(r.Products, productID)
	return nil
}

// Удаляет последний добавленный товар (LIFO), если приёмка открыта
func (r *Reception) RemoveLastProduct() (uuid.UUID, error) {
	if !r.IsOpen() {
		return uuid.Nil, errors.New("приёмка закрыта, нельзя удалять товары")
	}
	if len(r.Products) == 0 {
		return uuid.Nil, errors.New("нет товаров для удаления")
	}
	last := r.Products[len(r.Products)-1]
	r.Products = r.Products[:len(r.Products)-1]
	return last, nil
}

// Закрывает приёмку, если она открыта
func (r *Reception) Close() error {
	if !r.IsOpen() {
		return errors.New("приёмка уже закрыта")
	}
	r.Status = ReceptionClosed
	return nil
}
