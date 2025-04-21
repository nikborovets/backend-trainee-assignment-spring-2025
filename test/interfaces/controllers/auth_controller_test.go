package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/interfaces/controllers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockDummyLoginUC struct{ mock.Mock }

func (m *mockDummyLoginUC) Execute(ctx context.Context, role entities.UserRole) (string, error) {
	args := m.Called(ctx, role)
	return args.String(0), args.Error(1)
}

type mockRegisterUC struct{ mock.Mock }

func (m *mockRegisterUC) Execute(ctx context.Context, email, password string, role entities.UserRole) (entities.User, error) {
	args := m.Called(ctx, email, password, role)
	return args.Get(0).(entities.User), args.Error(1)
}

type mockLoginUC struct{ mock.Mock }

func (m *mockLoginUC) Execute(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func TestAuthController_DummyLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := new(mockDummyLoginUC)
	ctrl := controllers.NewAuthController(uc, nil, nil)
	r := gin.New()
	r.POST("/dummyLogin", ctrl.DummyLogin)

	// happy path
	body := `{"role":"moderator"}`
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx := req.Context()
	uc.On("Execute", ctx, entities.UserRoleModerator).Return("token123", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)
	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	require.Equal(t, "token123", resp["token"])
	uc.AssertExpectations(t)

	// ошибка usecase
	body = `{"role":"client"}`
	req = httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	ctx = req.Context()
	uc.On("Execute", ctx, entities.UserRoleClient).Return("", assert.AnError)
	r.ServeHTTP(w, req)
	require.Equal(t, 401, w.Code)
}

func TestAuthController_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := new(mockRegisterUC)
	ctrl := controllers.NewAuthController(nil, uc, nil)
	r := gin.New()
	r.POST("/register", ctrl.Register)

	user := entities.User{Email: "test@avito.ru", Role: entities.UserRoleModerator}
	body := `{"email":"test@avito.ru","password":"pass","role":"moderator"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx := req.Context()
	uc.On("Execute", ctx, "test@avito.ru", "pass", entities.UserRoleModerator).Return(user, nil)
	r.ServeHTTP(w, req)
	require.Equal(t, 201, w.Code)
	var resp entities.User
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	require.Equal(t, user.Email, resp.Email)
	uc.AssertExpectations(t)

	// ошибка usecase
	body = `{"email":"fail@avito.ru","password":"fail","role":"moderator"}`
	req = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	ctx = req.Context()
	uc.On("Execute", ctx, "fail@avito.ru", "fail", entities.UserRoleModerator).Return(entities.User{}, assert.AnError)
	r.ServeHTTP(w, req)
	require.Equal(t, 400, w.Code)
}

func TestAuthController_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := new(mockLoginUC)
	ctrl := controllers.NewAuthController(nil, nil, uc)
	r := gin.New()
	r.POST("/login", ctrl.Login)

	body := `{"email":"test@avito.ru","password":"pass"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx := req.Context()
	uc.On("Execute", ctx, "test@avito.ru", "pass").Return("token456", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)
	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	require.Equal(t, "token456", resp["token"])
	uc.AssertExpectations(t)

	// ошибка usecase
	body = `{"email":"fail@avito.ru","password":"fail"}`
	req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	ctx = req.Context()
	uc.On("Execute", ctx, "fail@avito.ru", "fail").Return("", assert.AnError)
	r.ServeHTTP(w, req)
	require.Equal(t, 401, w.Code)
}
