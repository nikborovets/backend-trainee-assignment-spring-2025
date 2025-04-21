package usecases

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"golang.org/x/crypto/bcrypt"
)

// UserRepositoryForRegister — интерфейс для создания пользователя
type UserRepositoryForRegister interface {
	Create(ctx context.Context, user entities.User, passwordHash string) (entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
}

// RegisterUseCase — интерактор для регистрации пользователя
type RegisterUseCase struct {
	repo UserRepositoryForRegister
}

func NewRegisterUseCase(repo UserRepositoryForRegister) *RegisterUseCase {
	return &RegisterUseCase{repo: repo}
}

// Execute регистрирует пользователя (email, пароль, роль)
func (uc *RegisterUseCase) Execute(ctx context.Context, email, password string, role entities.UserRole) (entities.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return entities.User{}, errors.New("email и пароль обязательны")
	}
	if !entities.ValidateUserRole(role) {
		return entities.User{}, errors.New("некорректная роль")
	}
	exists, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return entities.User{}, err
	}
	if exists != nil {
		return entities.User{}, errors.New("пользователь с таким email уже существует")
	}
	passwordHash, err := hashPassword(password)
	if err != nil {
		return entities.User{}, err
	}
	user := entities.User{
		ID:               uuid.New(),
		Email:            email,
		Role:             role,
		RegistrationDate: time.Now().UTC(),
	}
	return uc.repo.Create(ctx, user, passwordHash)
}

// hashPassword хэширует пароль через bcrypt
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
