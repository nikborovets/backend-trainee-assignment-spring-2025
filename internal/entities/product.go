package entities

import (
	"time"

	"github.com/google/uuid"
)

// Product — товар, принимаемый на ПВЗ
// id — UUID
// receptionId — UUID
// type — электроника/одежда/обувь
// receivedAt — дата и время приёма товара (момент добавления в систему)

type ProductType string

const (
	ProductElectronics ProductType = "электроника"
	ProductClothes     ProductType = "одежда"
	ProductShoes       ProductType = "обувь"
)

type Product struct {
	ID          uuid.UUID   `json:"id"`
	ReceptionID uuid.UUID   `json:"receptionId"`
	Type        ProductType `json:"type"`
	ReceivedAt  time.Time   `json:"receivedAt"`
}

// Проверяет, валиден ли тип товара
func ValidateProductType(t ProductType) bool {
	return t == ProductElectronics || t == ProductClothes || t == ProductShoes
}
