package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
)

type ReceptionController struct {
	CreateUC *usecases.CreateReceptionUseCase
}

func NewReceptionController(create *usecases.CreateReceptionUseCase) *ReceptionController {
	return &ReceptionController{CreateUC: create}
}

// POST /receptions {"pvzId": "..."}
func (c *ReceptionController) Create(ctx *gin.Context) {
	user := ctx.MustGet("user").(entities.User)
	var req struct {
		PVZID string `json:"pvzId"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	pvzID, err := uuid.Parse(req.PVZID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad pvzId"})
		return
	}
	rec, err := c.CreateUC.Execute(ctx, user, pvzID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, rec)
}
