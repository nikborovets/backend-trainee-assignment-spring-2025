package usecases

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/configs"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"golang.org/x/crypto/bcrypt"
)

// UserRepositoryForLogin — интерфейс для поиска пользователя по email с возвратом хэша пароля.
type UserRepositoryForLogin interface {
	GetByEmail(ctx context.Context, email string) (*entities.User, string, error)
}

// LoginUseCase — интерактор для логина по email+пароль, возвращает JWT.
type LoginUseCase struct {
	repo      UserRepositoryForLogin
	jwtSecret []byte
}

// NewLoginUseCase создаёт LoginUseCase с репозиторием и секретом JWT.
func NewLoginUseCase(repo UserRepositoryForLogin, cfg *configs.Config) *LoginUseCase {
	return &LoginUseCase{repo: repo, jwtSecret: []byte(cfg.JWTSecret)}
}

// Execute логинит пользователя по email+пароль, возвращает JWT-токен.
func (uc *LoginUseCase) Execute(ctx context.Context, email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return "", errors.New("email и пароль обязательны")
	}
	user, hash, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}
	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"role":  string(user.Role),
		"email": user.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtStr, err := token.SignedString(uc.jwtSecret)
	if err != nil {
		return "", err
	}
	return jwtStr, nil
}
