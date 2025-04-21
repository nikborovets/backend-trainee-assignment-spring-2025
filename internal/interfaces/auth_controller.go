package interfaces

import (
	"context"
)

// AuthController — интерфейс контроллера авторизации (см. .puml)
type AuthController interface {
	// DummyLogin выдаёт токен по роли (без пароля)
	DummyLogin(ctx context.Context, role string) (string, error)
	// Register регистрирует пользователя
	Register(ctx context.Context, req RegisterRequest) (UserDTO, error)
	// Login логинит пользователя по email+пароль
	Login(ctx context.Context, req LoginRequest) (string, error)
}
