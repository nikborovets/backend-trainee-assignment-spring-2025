package test

import (
	"testing"

	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

func TestValidateCity(t *testing.T) {
	// Arrange
	valid := []entities.City{entities.CityMoscow, entities.CitySPB, entities.CityKazan}
	invalid := []entities.City{"Воронеж", "Новосибирск", ""}

	// Act & Assert
	for _, c := range valid {
		if !entities.ValidateCity(c) {
			// Assert
			t.Errorf("city %s should be valid", c)
		}
	}
	for _, c := range invalid {
		if entities.ValidateCity(c) {
			// Assert
			t.Errorf("city %s should be invalid", c)
		}
	}
}

func TestValidateUserRole(t *testing.T) {
	// Arrange
	valid := []entities.UserRole{entities.UserRoleClient, entities.UserRoleModerator}
	invalid := []entities.UserRole{"pvz_staff", ""}

	// Act & Assert
	for _, r := range valid {
		if !entities.ValidateUserRole(r) {
			// Assert
			t.Errorf("%s should be valid role", r)
		}
	}
	for _, r := range invalid {
		if entities.ValidateUserRole(r) {
			// Assert
			t.Errorf("%s should be invalid role", r)
		}
	}
}

func TestValidateProductType(t *testing.T) {
	// Arrange
	valid := []entities.ProductType{entities.ProductElectronics, entities.ProductClothes, entities.ProductShoes}
	invalid := []entities.ProductType{"еда", ""}

	// Act & Assert
	for _, tpe := range valid {
		if !entities.ValidateProductType(tpe) {
			// Assert
			t.Errorf("%s should be valid type", tpe)
		}
	}
	for _, tpe := range invalid {
		if entities.ValidateProductType(tpe) {
			// Assert
			t.Errorf("%s should be invalid type", tpe)
		}
	}
}
