package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/configs"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepoForLogin struct {
	getByEmailFn func(ctx context.Context, email string) (*entities.User, string, error)
}

func (m *mockUserRepoForLogin) GetByEmail(ctx context.Context, email string) (*entities.User, string, error) {
	return m.getByEmailFn(ctx, email)
}

func TestLoginUseCase_Execute(t *testing.T) {
	// Arrange
	cfg := &configs.Config{JWTSecret: "testsecret"}
	user := &entities.User{ID: uuid.New(), Email: "test@avito.ru", Role: entities.UserRoleModerator, RegistrationDate: time.Now()}
	password := "password123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	repo := &mockUserRepoForLogin{
		getByEmailFn: func(ctx context.Context, email string) (*entities.User, string, error) {
			if email == user.Email {
				return user, string(hash), nil
			}
			return nil, "", nil
		},
	}
	uc := usecases.NewLoginUseCase(repo, cfg)
	ctx := context.Background()

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
		checkJWT bool
	}{
		{"ok", user.Email, password, false, true},
		{"wrong password", user.Email, "wrongpass", true, false},
		{"not found", "notfound@avito.ru", password, true, false},
		{"empty email", "", password, true, false},
		{"empty password", user.Email, "", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			token, err := uc.Execute(ctx, tt.email, tt.password)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, token)

			if tt.checkJWT {
				parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					return []byte(cfg.JWTSecret), nil
				})
				require.NoError(t, err)
				claims, ok := parsed.Claims.(jwt.MapClaims)
				require.True(t, ok)
				assert.Equal(t, user.Email, claims["email"])
				assert.Equal(t, string(user.Role), claims["role"])
			}
		})
	}
}
