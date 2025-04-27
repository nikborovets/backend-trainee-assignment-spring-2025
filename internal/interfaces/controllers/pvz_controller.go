package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/interfaces"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
)

type PVZController struct {
	CreateUC     usecases.CreatePVZUseCaseIface
	ListUC       usecases.ListPVZsUseCaseIface
	CloseUC      usecases.CloseReceptionUseCaseIface
	DeleteLastUC usecases.DeleteLastProductUseCaseIface
}

func NewPVZController(create usecases.CreatePVZUseCaseIface, list usecases.ListPVZsUseCaseIface, closeUC usecases.CloseReceptionUseCaseIface, delUC usecases.DeleteLastProductUseCaseIface) *PVZController {
	return &PVZController{
		CreateUC:     create,
		ListUC:       list,
		CloseUC:      closeUC,
		DeleteLastUC: delUC,
	}
}

// POST /pvz {"city": "Москва"}
func (c *PVZController) Create(ctx *gin.Context) {
	userVal, ok := ctx.Get("user")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	user := userVal.(entities.User)
	var req struct {
		City entities.City `json:"city"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	pvz, err := c.CreateUC.Execute(ctx.Request.Context(), user, req.City)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, pvz)
}

// GET /pvz?start=...&end=...&page=1&limit=10
func (c *PVZController) List(ctx *gin.Context) {
	userVal, ok := ctx.Get("user")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	user := userVal.(entities.User)
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
	pvzs, err := c.ListUC.Execute(ctx.Request.Context(), user, start, end, page, limit)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	// Используем новый DTO для правильного форматирования ответа
	var result []interfaces.PVZListResponseItem

	for _, pvz := range pvzs {
		// Создаем DTO для текущего ПВЗ
		pvzDTO := interfaces.ToPVZListItemDTO(pvz)

		// Для каждого PVZ — получить все приёмки и продукты
		receptions, _ := c.ListUC.GetReceptionsByPVZ(ctx.Request.Context(), pvz.ID) // предполагается, что usecase расширен
		var recs []struct {
			Reception interfaces.ReceptionListItemDTO `json:"reception"`
			Products  []interfaces.ProductDTO         `json:"products"`
		}

		for _, rec := range receptions {
			// Преобразуем Reception в ReceptionListItemDTO, исключая поле Products
			recDTO := interfaces.ToReceptionListItemDTO(rec)

			products, _ := c.ListUC.GetProductsByReception(ctx.Request.Context(), rec.ID)

			// Конвертируем доменные модели продуктов в DTO
			var productDTOs []interfaces.ProductDTO
			for _, product := range products {
				productDTOs = append(productDTOs, interfaces.ToProductDTO(product))
			}

			recs = append(recs, struct {
				Reception interfaces.ReceptionListItemDTO `json:"reception"`
				Products  []interfaces.ProductDTO         `json:"products"`
			}{Reception: recDTO, Products: productDTOs})
		}

		result = append(result, interfaces.PVZListResponseItem{
			PVZ:        pvzDTO,
			Receptions: recs,
		})
	}
	ctx.JSON(http.StatusOK, result)
}

// POST /pvz/:pvzId/close_last_reception
func (c *PVZController) CloseLastReception(ctx *gin.Context) {
	userVal, ok := ctx.Get("user")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	user := userVal.(entities.User)
	pvzID, err := uuid.Parse(ctx.Param("pvzId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad pvzId"})
		return
	}
	rec, err := c.CloseUC.Execute(ctx.Request.Context(), user, pvzID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, rec)
}

// POST /pvz/:pvzId/delete_last_product
func (c *PVZController) DeleteLastProduct(ctx *gin.Context) {
	userVal, ok := ctx.Get("user")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	user := userVal.(entities.User)
	pvzID, err := uuid.Parse(ctx.Param("pvzId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad pvzId"})
		return
	}
	err = c.DeleteLastUC.Execute(ctx.Request.Context(), user, pvzID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}
