package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/configs"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

// DummyLoginUseCase выдаёт токен по роли без проверки email/пароля.
type DummyLoginUseCase struct {
	jwtSecret []byte
}

// NewDummyLoginUseCase создаёт DummyLoginUseCase с секретом для JWT из конфига.
func NewDummyLoginUseCase(cfg *configs.Config) *DummyLoginUseCase {
	return &DummyLoginUseCase{jwtSecret: []byte(cfg.JWTSecret)}
}

// Execute выдаёт JWT-токен для заданной роли. Возвращает ошибку, если роль невалидна.
func (uc *DummyLoginUseCase) Execute(ctx context.Context, role entities.UserRole) (string, error) {
	if !entities.ValidateUserRole(role) {
		return "", ErrInvalidRole
	}

	claims := jwt.MapClaims{
		"sub":  uuid.NewString(),
		"role": string(role),
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtStr, err := token.SignedString(uc.jwtSecret)
	if err != nil {
		return "", err
	}
	return jwtStr, nil
}

// ErrInvalidRole возвращается, если роль невалидна.
var ErrInvalidRole = errors.New("invalid user role")
