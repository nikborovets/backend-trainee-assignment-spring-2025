package interfaces

import (
	"time"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

// UserDTO — DTO пользователя для API
// id, email, role
// swagger:model UserDTO
// example: {"id":"uuid","email":"user@avito.ru","role":"client"}
type UserDTO struct {
	ID    uuid.UUID         `json:"id"`
	Email string            `json:"email"`
	Role  entities.UserRole `json:"role"`
}

// PVZDTO — DTO ПВЗ для API
// id, registrationDate, city
type PVZDTO struct {
	ID               uuid.UUID     `json:"id"`
	RegistrationDate time.Time     `json:"registrationDate"`
	City             entities.City `json:"city"`
}

// FullPVZDTO — ПВЗ + все приёмки с товарами
type FullPVZDTO struct {
	PVZ        PVZDTO                     `json:"pvz"`
	Receptions []ReceptionWithProductsDTO `json:"receptions"`
}

// ReceptionDTO — DTO приёмки
type ReceptionDTO struct {
	ID       uuid.UUID                `json:"id"`
	DateTime time.Time                `json:"dateTime"`
	Status   entities.ReceptionStatus `json:"status"`
	PVZID    uuid.UUID                `json:"pvzId"`
}

// ReceptionListItemDTO - структура для возврата Reception в API ответе без поля products
type ReceptionListItemDTO struct {
	ID       uuid.UUID                `json:"id"`
	PVZID    uuid.UUID                `json:"pvzId"`
	Status   entities.ReceptionStatus `json:"status"`
	DateTime time.Time                `json:"dateTime"`
	// Поле Products здесь специально отсутствует
}

// ToReceptionListItemDTO конвертирует Reception в DTO для API ответа
func ToReceptionListItemDTO(reception entities.Reception) ReceptionListItemDTO {
	return ReceptionListItemDTO{
		ID:       reception.ID,
		PVZID:    reception.PVZID,
		Status:   reception.Status,
		DateTime: reception.DateTime,
	}
}

// ProductDTO — DTO товара
type ProductDTO struct {
	ID          uuid.UUID            `json:"id"`
	DateTime    time.Time            `json:"dateTime"`
	Type        entities.ProductType `json:"type"`
	ReceptionID uuid.UUID            `json:"receptionId"`
}

// ToProductDTO преобразует доменную модель Product в DTO для API
func ToProductDTO(product entities.Product) ProductDTO {
	return ProductDTO{
		ID:          product.ID,
		DateTime:    product.DateTime,
		Type:        product.Type,
		ReceptionID: product.ReceptionID,
	}
}

// ReceptionWithProductsDTO — приёмка + товары
type ReceptionWithProductsDTO struct {
	Reception ReceptionDTO `json:"reception"`
	Products  []ProductDTO `json:"products"`
}

// RegisterRequest — DTO для регистрации пользователя
type RegisterRequest struct {
	Email    string            `json:"email"`
	Password string            `json:"password"`
	Role     entities.UserRole `json:"role"`
}

// LoginRequest — DTO для логина
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ListParams — параметры фильтрации/пагинации для /pvz
type ListParams struct {
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
	Page      int        `json:"page"`
	Limit     int        `json:"limit"`
}

// AddProductRequest — DTO для добавления товара
type AddProductRequest struct {
	PVZID uuid.UUID            `json:"pvzId"`
	Type  entities.ProductType `json:"type"`
}

// PVZListItemDTO - структура для возврата ПВЗ в API ответе /pvz GET
type PVZListItemDTO struct {
	ID               uuid.UUID     `json:"id"`
	RegistrationDate time.Time     `json:"registrationDate"`
	City             entities.City `json:"city"`
	// Поле Receptions здесь специально отсутствует
}

// ToPVZListItemDTO конвертирует PVZ в DTO для API ответа
func ToPVZListItemDTO(pvz entities.PVZ) PVZListItemDTO {
	return PVZListItemDTO{
		ID:               pvz.ID,
		RegistrationDate: pvz.RegistrationDate,
		City:             pvz.City,
	}
}

// PVZListResponseItem - структура для элемента в ответе списка ПВЗ
type PVZListResponseItem struct {
	PVZ        PVZListItemDTO `json:"pvz"`
	Receptions []struct {
		Reception ReceptionListItemDTO `json:"reception"` // Используем ReceptionListItemDTO вместо entities.Reception
		Products  []ProductDTO         `json:"products"`
	} `json:"receptions"`
}
