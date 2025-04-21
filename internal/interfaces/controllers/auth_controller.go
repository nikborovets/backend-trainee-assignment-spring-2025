package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
)

type AuthController struct {
	DummyLoginUC usecases.DummyLoginUseCaseIface
	RegisterUC   usecases.RegisterUseCaseIface
	LoginUC      usecases.LoginUseCaseIface
}

func NewAuthController(dummy usecases.DummyLoginUseCaseIface, reg usecases.RegisterUseCaseIface, login usecases.LoginUseCaseIface) *AuthController {
	return &AuthController{
		DummyLoginUC: dummy,
		RegisterUC:   reg,
		LoginUC:      login,
	}
}

// POST /dummyLogin {"role": "moderator"}
func (c *AuthController) DummyLogin(ctx *gin.Context) {
	var req struct {
		Role entities.UserRole `json:"role"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	token, err := c.DummyLoginUC.Execute(ctx.Request.Context(), req.Role)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// POST /register {"email":..., "password":..., "role":...}
func (c *AuthController) Register(ctx *gin.Context) {
	var req struct {
		Email    string            `json:"email"`
		Password string            `json:"password"`
		Role     entities.UserRole `json:"role"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	user, err := c.RegisterUC.Execute(ctx.Request.Context(), req.Email, req.Password, req.Role)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, user)
}

// POST /login {"email":..., "password":...}
func (c *AuthController) Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	token, err := c.LoginUC.Execute(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
