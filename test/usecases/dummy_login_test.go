package usecases_test

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/configs"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
	"github.com/stretchr/testify/require"
)

func TestDummyLoginUseCase_Execute(t *testing.T) {
	// Arrange
	cfg := &configs.Config{JWTSecret: "testsecret"}
	uc := usecases.NewDummyLoginUseCase(cfg)
	ctx := context.Background()

	tests := []struct {
		name    string
		role    entities.UserRole
		wantErr bool
	}{
		{"valid_client", entities.UserRoleClient, false},
		{"valid_moderator", entities.UserRoleModerator, false},
		{"valid_pvz_staff", entities.UserRolePVZStaff, false},
		{"invalid_role", "hacker", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			tokenStr, err := uc.Execute(ctx, tt.role)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, tokenStr)

			// Проверяем, что токен валиден и роль совпадает
			parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			})
			require.NoError(t, err)
			claims, ok := parsed.Claims.(jwt.MapClaims)
			require.True(t, ok)
			roleClaim, ok := claims["role"].(string)
			require.True(t, ok)
			if !tt.wantErr {
				require.Equal(t, string(tt.role), roleClaim)
			}
			// Проверяем exp/iat
			_, hasExp := claims["exp"]
			_, hasIat := claims["iat"]
			require.True(t, hasExp)
			require.True(t, hasIat)
		})
	}
}
