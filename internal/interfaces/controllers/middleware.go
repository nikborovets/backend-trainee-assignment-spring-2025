package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h := ctx.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing or invalid Authorization header"})
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid claims"})
			return
		}
		user := entities.User{
			Email: claims["email"].(string),
			Role:  entities.UserRole(claims["role"].(string)),
		}
		ctx.Set("user", user)
		ctx.Next()
	}
}
