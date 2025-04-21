package test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

func TestReceptionIsOpen(t *testing.T) {
	r := entities.Reception{Status: entities.ReceptionInProgress}
	if !r.IsOpen() {
		t.Error("Reception should be open")
	}
	r.Status = entities.ReceptionClosed
	if r.IsOpen() {
		t.Error("Reception should be closed")
	}
}

func TestReceptionAddProduct(t *testing.T) {
	r := entities.Reception{Status: entities.ReceptionInProgress}
	pid := uuid.New()
	if err := r.AddProduct(pid); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(r.Products) != 1 || r.Products[0] != pid {
		t.Error("product not added correctly")
	}
	r.Status = entities.ReceptionClosed
	if err := r.AddProduct(uuid.New()); err == nil {
		t.Error("should not add product to closed reception")
	}
}

func TestReceptionRemoveLastProduct(t *testing.T) {
	r := entities.Reception{Status: entities.ReceptionInProgress}
	p1, p2 := uuid.New(), uuid.New()
	r.Products = []uuid.UUID{p1, p2}
	removed, err := r.RemoveLastProduct()
	if err != nil || removed != p2 {
		t.Error("should remove last product (LIFO)")
	}
	if len(r.Products) != 1 || r.Products[0] != p1 {
		t.Error("product list not updated after remove")
	}
	r.Status = entities.ReceptionClosed
	if _, err := r.RemoveLastProduct(); err == nil {
		t.Error("should not remove from closed reception")
	}
}

func TestReceptionClose(t *testing.T) {
	r := entities.Reception{Status: entities.ReceptionInProgress}
	if err := r.Close(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if r.Status != entities.ReceptionClosed {
		t.Error("reception should be closed after Close()")
	}
	if err := r.Close(); err == nil {
		t.Error("should not close already closed reception")
	}
}
