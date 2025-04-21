package entities

import (
	"time"

	"github.com/google/uuid"
)

// PVZ — пункт выдачи заказов
// Может быть только в городах: Москва, Санкт-Петербург, Казань
// registrationDate — дата регистрации
// receptions — список приёмок (UUID)
type PVZ struct {
	ID               uuid.UUID   `json:"id"`
	RegistrationDate time.Time   `json:"registrationDate"`
	City             City        `json:"city"`
	Receptions       []uuid.UUID `json:"receptions"`
}

type City string

const (
	CityMoscow City = "Москва"
	CitySPB    City = "Санкт-Петербург"
	CityKazan  City = "Казань"
)

var AllowedPVZCities = []City{CityMoscow, CitySPB, CityKazan}

// Проверяет, разрешён ли город для PVZ
func ValidateCity(city City) bool {
	for _, allowed := range AllowedPVZCities {
		if city == allowed {
			return true
		}
	}
	return false
}
