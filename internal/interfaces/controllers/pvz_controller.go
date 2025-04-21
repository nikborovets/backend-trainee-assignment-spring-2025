package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
)

type PVZController struct {
	CreateUC     *usecases.CreatePVZUseCase
	ListUC       *usecases.ListPVZsUseCase
	CloseUC      *usecases.CloseReceptionUseCase
	DeleteLastUC *usecases.DeleteLastProductUseCase
}

func NewPVZController(create *usecases.CreatePVZUseCase, list *usecases.ListPVZsUseCase, closeUC *usecases.CloseReceptionUseCase, delUC *usecases.DeleteLastProductUseCase) *PVZController {
	return &PVZController{
		CreateUC:     create,
		ListUC:       list,
		CloseUC:      closeUC,
		DeleteLastUC: delUC,
	}
}

// POST /pvz {"city": "Москва"}
func (c *PVZController) Create(ctx *gin.Context) {
	user := ctx.MustGet("user").(entities.User)
	var req struct {
		City entities.City `json:"city"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	pvz, err := c.CreateUC.Execute(ctx, user, req.City)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, pvz)
}

// GET /pvz?start=...&end=...&page=1&limit=10
func (c *PVZController) List(ctx *gin.Context) {
	user := ctx.MustGet("user").(entities.User)
	var start, end *time.Time
	if s := ctx.Query("start"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err == nil {
			start = &t
		}
	}
	if e := ctx.Query("end"); e != "" {
		t, err := time.Parse(time.RFC3339, e)
		if err == nil {
			end = &t
		}
	}
	page := 1
	limit := 10
	if p := ctx.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := ctx.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	// --- агрегирующий usecase ---
	pvzs, err := c.ListUC.Execute(ctx, user, start, end, page, limit)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}
	var result []struct {
		PVZ        entities.PVZ `json:"pvz"`
		Receptions []struct {
			Reception entities.Reception `json:"reception"`
			Products  []entities.Product `json:"products"`
		} `json:"receptions"`
	}
	for _, pvz := range pvzs {
		// Для каждого PVZ — получить все приёмки и продукты
		receptions, _ := c.ListUC.GetReceptionsByPVZ(ctx, pvz.ID) // предполагается, что usecase расширен
		var recs []struct {
			Reception entities.Reception `json:"reception"`
			Products  []entities.Product `json:"products"`
		}
		for _, rec := range receptions {
			products, _ := c.ListUC.GetProductsByReception(ctx, rec.ID)
			recs = append(recs, struct {
				Reception entities.Reception `json:"reception"`
				Products  []entities.Product `json:"products"`
			}{Reception: rec, Products: products})
		}
		result = append(result, struct {
			PVZ        entities.PVZ `json:"pvz"`
			Receptions []struct {
				Reception entities.Reception `json:"reception"`
				Products  []entities.Product `json:"products"`
			} `json:"receptions"`
		}{PVZ: pvz, Receptions: recs})
	}
	ctx.JSON(http.StatusOK, result)
}

// POST /pvz/:pvzId/close_last_reception
func (c *PVZController) CloseLastReception(ctx *gin.Context) {
	user := ctx.MustGet("user").(entities.User)
	pvzID, err := uuid.Parse(ctx.Param("pvzId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad pvzId"})
		return
	}
	rec, err := c.CloseUC.Execute(ctx, user, pvzID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, rec)
}

// POST /pvz/:pvzId/delete_last_product
func (c *PVZController) DeleteLastProduct(ctx *gin.Context) {
	user := ctx.MustGet("user").(entities.User)
	pvzID, err := uuid.Parse(ctx.Param("pvzId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad pvzId"})
		return
	}
	err = c.DeleteLastUC.Execute(ctx, user, pvzID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}
