package entities_test

import (
	"testing"

	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/stretchr/testify/require"
)

func TestValidateCity(t *testing.T) {
	// Arrange
	valid := []entities.City{entities.CityMoscow, entities.CitySPB, entities.CityKazan}
	invalid := []entities.City{"Тверь", "Воронеж", ""}

	// Act & Assert
	for _, c := range valid {
		require.True(t, entities.ValidateCity(c), "city %s should be valid", c)
	}
	for _, c := range invalid {
		require.False(t, entities.ValidateCity(c), "city %s should be invalid", c)
	}
}

func TestValidateUserRole(t *testing.T) {
	// Arrange
	valid := []entities.UserRole{entities.UserRoleClient, entities.UserRoleModerator, entities.UserRolePVZStaff}
	invalid := []entities.UserRole{"admin", "", "hacker"}

	// Act & Assert
	for _, r := range valid {
		require.True(t, entities.ValidateUserRole(r), "%s should be valid role", r)
	}
	for _, r := range invalid {
		require.False(t, entities.ValidateUserRole(r), "%s should be invalid role", r)
	}
}

func TestValidateProductType(t *testing.T) {
	// Arrange
	valid := []entities.ProductType{entities.ProductElectronics, entities.ProductClothes, entities.ProductShoes}
	invalid := []entities.ProductType{"еда", "", "машина"}

	// Act & Assert
	for _, tpe := range valid {
		require.True(t, entities.ValidateProductType(tpe), "%s should be valid type", tpe)
	}
	for _, tpe := range invalid {
		require.False(t, entities.ValidateProductType(tpe), "%s should be invalid type", tpe)
	}
}
