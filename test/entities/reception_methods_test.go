package entities_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReceptionIsOpen(t *testing.T) {
	// Arrange
	r := entities.Reception{Status: entities.ReceptionInProgress}

	// Act & Assert
	require.True(t, r.IsOpen(), "Reception should be open")
	r.Status = entities.ReceptionClosed
	require.False(t, r.IsOpen(), "Reception should be closed")
}

func TestReceptionAddProduct(t *testing.T) {
	// Arrange
	r := entities.Reception{Status: entities.ReceptionInProgress}
	pid := uuid.New()

	// Act
	err := r.AddProduct(pid)

	// Assert
	require.NoError(t, err, "unexpected error")
	require.Len(t, r.Products, 1)
	require.Equal(t, pid, r.Products[0], "product not added correctly")

	// Arrange (закрываем приёмку)
	r.Status = entities.ReceptionClosed

	// Act
	err = r.AddProduct(uuid.New())

	// Assert
	assert.Error(t, err, "should not add product to closed reception")
}

func TestReceptionRemoveLastProduct(t *testing.T) {
	// Arrange
	r := entities.Reception{Status: entities.ReceptionInProgress}
	p1, p2 := uuid.New(), uuid.New()
	r.Products = []uuid.UUID{p1, p2}

	// Act
	removed, err := r.RemoveLastProduct()

	// Assert
	require.NoError(t, err)
	require.Equal(t, p2, removed, "should remove last product (LIFO)")
	require.Len(t, r.Products, 1)
	require.Equal(t, p1, r.Products[0], "product list not updated after remove")

	// Arrange (закрываем приёмку)
	r.Status = entities.ReceptionClosed

	// Act
	_, err = r.RemoveLastProduct()

	// Assert
	assert.Error(t, err, "should not remove from closed reception")
}

func TestReceptionClose(t *testing.T) {
	// Arrange
	r := entities.Reception{Status: entities.ReceptionInProgress}

	// Act
	err := r.Close()

	// Assert
	require.NoError(t, err, "unexpected error")
	require.Equal(t, entities.ReceptionClosed, r.Status, "reception should be closed after Close()")

	// Act (пытаемся закрыть ещё раз)
	err = r.Close()

	// Assert
	assert.Error(t, err, "should not close already closed reception")
}
