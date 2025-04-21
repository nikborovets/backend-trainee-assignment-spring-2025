package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
)

type ProductController struct {
	AddUC *usecases.AddProductUseCase
}

func NewProductController(add *usecases.AddProductUseCase) *ProductController {
	return &ProductController{AddUC: add}
}

// POST /products {"pvzId": "...", "type": "электроника"}
func (c *ProductController) Add(ctx *gin.Context) {
	user := ctx.MustGet("user").(entities.User)
	var req struct {
		PVZID string               `json:"pvzId"`
		Type  entities.ProductType `json:"type"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	pvzID, err := uuid.Parse(req.PVZID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad pvzId"})
		return
	}
	product, err := c.AddUC.Execute(ctx, user, pvzID, req.Type)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, product)
}
