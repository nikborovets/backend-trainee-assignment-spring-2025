package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/interfaces/controllers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockCreatePVZUC struct{ mock.Mock }

func (m *mockCreatePVZUC) Execute(ctx context.Context, user entities.User, city entities.City) (entities.PVZ, error) {
	args := m.Called(ctx, user, city)
	return args.Get(0).(entities.PVZ), args.Error(1)
}

type mockListPVZsUC struct{ mock.Mock }

func (m *mockListPVZsUC) Execute(ctx context.Context, user entities.User, start, end *time.Time, page, limit int) ([]entities.PVZ, error) {
	args := m.Called(ctx, user, start, end, page, limit)
	return args.Get(0).([]entities.PVZ), args.Error(1)
}
func (m *mockListPVZsUC) GetReceptionsByPVZ(ctx context.Context, pvzID uuid.UUID) ([]entities.Reception, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).([]entities.Reception), args.Error(1)
}
func (m *mockListPVZsUC) GetProductsByReception(ctx context.Context, recID uuid.UUID) ([]entities.Product, error) {
	args := m.Called(ctx, recID)
	return args.Get(0).([]entities.Product), args.Error(1)
}

type mockCloseReceptionUC struct{ mock.Mock }

func (m *mockCloseReceptionUC) Execute(ctx context.Context, user entities.User, pvzID uuid.UUID) (entities.Reception, error) {
	args := m.Called(ctx, user, pvzID)
	return args.Get(0).(entities.Reception), args.Error(1)
}

type mockDeleteLastProductUC struct{ mock.Mock }

func (m *mockDeleteLastProductUC) Execute(ctx context.Context, user entities.User, pvzID uuid.UUID) error {
	args := m.Called(ctx, user, pvzID)
	return args.Error(0)
}

func TestPVZController_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := new(mockCreatePVZUC)
	ctrl := controllers.NewPVZController(uc, nil, nil, nil)
	r := gin.New()
	r.POST("/pvz", func(ctx *gin.Context) {
		ctx.Set("user", entities.User{Role: entities.UserRoleModerator})
		ctrl.Create(ctx)
	})

	t.Run("happy path", func(t *testing.T) {
		city := entities.City("Москва")
		user := entities.User{Role: entities.UserRoleModerator}
		pvz := entities.PVZ{ID: uuid.New(), City: city}
		uc.On("Execute", mock.MatchedBy(func(ctx context.Context) bool { return true }), user, city).Return(pvz, nil)
		body := `{"city":"Москва"}`
		req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, 201, w.Code)
		var resp entities.PVZ
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		require.Equal(t, city, resp.City)
		uc.AssertExpectations(t)
	})

	t.Run("ошибка usecase", func(t *testing.T) {
		city := entities.City("Москва")
		user := entities.User{Role: entities.UserRoleModerator}
		uc := new(mockCreatePVZUC)
		ctrl := controllers.NewPVZController(uc, nil, nil, nil)
		r := gin.New()
		r.POST("/pvz", func(ctx *gin.Context) {
			ctx.Set("user", user)
			ctrl.Create(ctx)
		})
		uc.On("Execute", mock.MatchedBy(func(ctx context.Context) bool { return true }), user, city).Return(entities.PVZ{}, assert.AnError)
		body := `{"city":"Москва"}`
		req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, 403, w.Code)
	})
}

func TestPVZController_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := new(mockListPVZsUC)
	ctrl := controllers.NewPVZController(nil, uc, nil, nil)
	r := gin.New()
	r.GET("/pvz", func(ctx *gin.Context) {
		ctx.Set("user", entities.User{Role: entities.UserRoleModerator})
		ctrl.List(ctx)
	})

	user := entities.User{Role: entities.UserRoleModerator}
	pvzID := uuid.New()
	pvz := entities.PVZ{ID: pvzID, City: "Москва"}
	recID := uuid.New()
	rec := entities.Reception{ID: recID}
	product := entities.Product{ID: uuid.New()}
	uc.On("Execute", mock.MatchedBy(func(ctx context.Context) bool { return true }), user, mock.Anything, mock.Anything, 1, 10).Return([]entities.PVZ{pvz}, nil)
	uc.On("GetReceptionsByPVZ", mock.MatchedBy(func(ctx context.Context) bool { return true }), pvzID).Return([]entities.Reception{rec}, nil)
	uc.On("GetProductsByReception", mock.MatchedBy(func(ctx context.Context) bool { return true }), recID).Return([]entities.Product{product}, nil)

	req := httptest.NewRequest(http.MethodGet, "/pvz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)
	var resp []map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 1)
	uc.AssertExpectations(t)
}

func TestPVZController_CloseLastReception(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("happy path", func(t *testing.T) {
		uc := new(mockCloseReceptionUC)
		ctrl := controllers.NewPVZController(nil, nil, uc, nil)
		r := gin.New()
		r.POST("/pvz/:pvzId/close_last_reception", func(ctx *gin.Context) {
			ctx.Set("user", entities.User{Role: entities.UserRolePVZStaff})
			ctrl.CloseLastReception(ctx)
		})
		user := entities.User{Role: entities.UserRolePVZStaff}
		pvzID := uuid.New()
		rec := entities.Reception{ID: uuid.New()}
		uc.On("Execute", mock.MatchedBy(func(ctx context.Context) bool { return true }), user, pvzID).Return(rec, nil)
		url := "/pvz/" + pvzID.String() + "/close_last_reception"
		req := httptest.NewRequest(http.MethodPost, url, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, 200, w.Code)
		var resp entities.Reception
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		require.Equal(t, rec.ID, resp.ID)
		uc.AssertExpectations(t)
	})

	t.Run("ошибка usecase", func(t *testing.T) {
		uc := new(mockCloseReceptionUC)
		ctrl := controllers.NewPVZController(nil, nil, uc, nil)
		r := gin.New()
		r.POST("/pvz/:pvzId/close_last_reception", func(ctx *gin.Context) {
			ctx.Set("user", entities.User{Role: entities.UserRolePVZStaff})
			ctrl.CloseLastReception(ctx)
		})
		user := entities.User{Role: entities.UserRolePVZStaff}
		pvzID := uuid.New()
		uc.On("Execute", mock.MatchedBy(func(ctx context.Context) bool { return true }), user, pvzID).Return(entities.Reception{}, assert.AnError)
		url := "/pvz/" + pvzID.String() + "/close_last_reception"
		req := httptest.NewRequest(http.MethodPost, url, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, 400, w.Code)
	})
}

func TestPVZController_DeleteLastProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("happy path", func(t *testing.T) {
		uc := new(mockDeleteLastProductUC)
		ctrl := controllers.NewPVZController(nil, nil, nil, uc)
		r := gin.New()
		r.POST("/pvz/:pvzId/delete_last_product", func(ctx *gin.Context) {
			ctx.Set("user", entities.User{Role: entities.UserRolePVZStaff})
			ctrl.DeleteLastProduct(ctx)
		})
		user := entities.User{Role: entities.UserRolePVZStaff}
		pvzID := uuid.New()
		uc.On("Execute", mock.MatchedBy(func(ctx context.Context) bool { return true }), user, pvzID).Return(nil)
		url := "/pvz/" + pvzID.String() + "/delete_last_product"
		req := httptest.NewRequest(http.MethodPost, url, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, 200, w.Code)
		var resp map[string]bool
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp["ok"])
		uc.AssertExpectations(t)
	})

	t.Run("ошибка usecase", func(t *testing.T) {
		uc := new(mockDeleteLastProductUC)
		ctrl := controllers.NewPVZController(nil, nil, nil, uc)
		r := gin.New()
		r.POST("/pvz/:pvzId/delete_last_product", func(ctx *gin.Context) {
			ctx.Set("user", entities.User{Role: entities.UserRolePVZStaff})
			ctrl.DeleteLastProduct(ctx)
		})
		user := entities.User{Role: entities.UserRolePVZStaff}
		pvzID := uuid.New()
		uc.On("Execute", mock.MatchedBy(func(ctx context.Context) bool { return true }), user, pvzID).Return(assert.AnError)
		url := "/pvz/" + pvzID.String() + "/delete_last_product"
		req := httptest.NewRequest(http.MethodPost, url, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, 400, w.Code)
	})
}
