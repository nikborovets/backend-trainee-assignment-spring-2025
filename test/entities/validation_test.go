package test

import (
	"testing"

	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

func TestValidateCity(t *testing.T) {
	valid := []entities.City{entities.CityMoscow, entities.CitySPB, entities.CityKazan}
	invalid := []entities.City{"Воронеж", "Новосибирск", ""}
	for _, c := range valid {
		if !entities.ValidateCity(c) {
			t.Errorf("city %s should be valid", c)
		}
	}
	for _, c := range invalid {
		if entities.ValidateCity(c) {
			t.Errorf("city %s should be invalid", c)
		}
	}
}

func TestValidateUserRole(t *testing.T) {
	if !entities.ValidateUserRole(entities.UserRoleClient) {
		t.Error("client should be valid role")
	}
	if !entities.ValidateUserRole(entities.UserRoleModerator) {
		t.Error("moderator should be valid role")
	}
	if entities.ValidateUserRole("pvz_staff") {
		t.Error("pvz_staff should be invalid role")
	}
	if entities.ValidateUserRole("") {
		t.Error("empty role should be invalid")
	}
}

func TestValidateProductType(t *testing.T) {
	if !entities.ValidateProductType(entities.ProductElectronics) {
		t.Error("электроника should be valid type")
	}
	if !entities.ValidateProductType(entities.ProductClothes) {
		t.Error("одежда should be valid type")
	}
	if !entities.ValidateProductType(entities.ProductShoes) {
		t.Error("обувь should be valid type")
	}
	if entities.ValidateProductType("еда") {
		t.Error("еда should be invalid type")
	}
	if entities.ValidateProductType("") {
		t.Error("empty type should be invalid")
	}
}
