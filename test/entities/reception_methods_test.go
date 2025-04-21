package test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

func TestReceptionIsOpen(t *testing.T) {
	// Arrange
	r := entities.Reception{Status: entities.ReceptionInProgress}

	// Act & Assert
	if !r.IsOpen() {
		t.Error("Reception should be open")
	}
	r.Status = entities.ReceptionClosed
	if r.IsOpen() {
		t.Error("Reception should be closed")
	}
}

func TestReceptionAddProduct(t *testing.T) {
	// Arrange
	r := entities.Reception{Status: entities.ReceptionInProgress}
	pid := uuid.New()

	// Act
	err := r.AddProduct(pid)

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(r.Products) != 1 || r.Products[0] != pid {
		t.Error("product not added correctly")
	}

	// Arrange (закрываем приёмку)
	r.Status = entities.ReceptionClosed

	// Act
	err = r.AddProduct(uuid.New())

	// Assert
	if err == nil {
		t.Error("should not add product to closed reception")
	}
}

func TestReceptionRemoveLastProduct(t *testing.T) {
	// Arrange
	r := entities.Reception{Status: entities.ReceptionInProgress}
	p1, p2 := uuid.New(), uuid.New()
	r.Products = []uuid.UUID{p1, p2}

	// Act
	removed, err := r.RemoveLastProduct()

	// Assert
	if err != nil || removed != p2 {
		t.Error("should remove last product (LIFO)")
	}
	if len(r.Products) != 1 || r.Products[0] != p1 {
		t.Error("product list not updated after remove")
	}

	// Arrange (закрываем приёмку)
	r.Status = entities.ReceptionClosed

	// Act
	_, err = r.RemoveLastProduct()

	// Assert
	if err == nil {
		t.Error("should not remove from closed reception")
	}
}

func TestReceptionClose(t *testing.T) {
	// Arrange
	r := entities.Reception{Status: entities.ReceptionInProgress}

	// Act
	err := r.Close()

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if r.Status != entities.ReceptionClosed {
		t.Error("reception should be closed after Close()")
	}

	// Act (пытаемся закрыть ещё раз)
	err = r.Close()

	// Assert
	if err == nil {
		t.Error("should not close already closed reception")
	}
}
