package entities

import (
	"time"

	"github.com/google/uuid"
)

// User — пользователь системы (client, moderator)
type User struct {
	ID               uuid.UUID `json:"id"`
	Email            string    `json:"email"`
	Role             UserRole  `json:"role"`
	RegistrationDate time.Time `json:"registrationDate"`
}

// UserRole — роль пользователя.
type UserRole string

const (
	UserRoleClient    UserRole = "client"
	UserRoleModerator UserRole = "moderator"
)

// Проверяет, валидна ли роль пользователя
func ValidateUserRole(role UserRole) bool {
	return role == UserRoleClient || role == UserRoleModerator
}
