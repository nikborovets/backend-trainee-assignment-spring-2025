package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/configs"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/infrastructure/repositories"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/interfaces/controllers"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/usecases"
	"github.com/stretchr/testify/require"
)

func init() {
	// Загружаем переменные окружения из .env файла
	_ = godotenv.Load("../../.env")

	// Вывод значения для отладки
	dsn := os.Getenv("TEST_PG_DSN")
	if dsn != "" {
		println("TEST_PG_DSN loaded in api_integration_test")
	}
}

// Адаптеры для репозиториев
type productRepoAdapter struct {
	*repositories.PGProductRepository
}

// DeleteLast только для интерфейса interfaces.ProductRepository
func (r *productRepoAdapter) DeleteLast(ctx context.Context, receptionID uuid.UUID) error {
	_, err := r.PGProductRepository.DeleteLast(ctx, receptionID)
	return err
}

type productRepoForDelete struct {
	*repositories.PGProductRepository
}

// DeleteLast для интерфейса usecases.ProductRepositoryForDelete
func (r *productRepoForDelete) DeleteLast(ctx context.Context, receptionID uuid.UUID) (*entities.Product, error) {
	return r.PGProductRepository.DeleteLast(ctx, receptionID)
}

type receptionRepoAdapter struct {
	*repositories.PGReceptionRepository
}

func (r *receptionRepoAdapter) CloseLast(ctx context.Context, pvzID uuid.UUID) (entities.Reception, error) {
	err := r.PGReceptionRepository.CloseLast(ctx, pvzID, time.Now().UTC())
	if err != nil {
		return entities.Reception{}, err
	}
	rec, err := r.GetActive(ctx, pvzID)
	if err != nil {
		return entities.Reception{}, err
	}
	if rec == nil {
		return entities.Reception{}, sql.ErrNoRows
	}
	return *rec, nil
}

func setupTestServer(t *testing.T) (*gin.Engine, *sql.DB) {
	db := setupTestDB(t)

	// Инициализация репозиториев
	pvzRepo := repositories.NewPGPVZRepository(db)
	receptionRepo := &receptionRepoAdapter{repositories.NewPGReceptionRepository(db)}
	productRepo := &productRepoAdapter{repositories.NewPGProductRepository(db)}
	productRepoDelete := &productRepoForDelete{repositories.NewPGProductRepository(db)}

	// Инициализация use cases
	createPVZUC := usecases.NewCreatePVZUseCase(pvzRepo)
	listPVZsUC := usecases.NewListPVZsUseCase(pvzRepo, receptionRepo, productRepo)
	closeReceptionUC := usecases.NewCloseReceptionUseCase(receptionRepo)
	deleteLastProductUC := usecases.NewDeleteLastProductUseCase(productRepoDelete, receptionRepo)
	createReceptionUC := usecases.NewCreateReceptionUseCase(receptionRepo)
	addProductUC := usecases.NewAddProductUseCase(productRepo, receptionRepo)
	dummyLoginUC := usecases.NewDummyLoginUseCase(&configs.Config{JWTSecret: "test_secret"})

	// Инициализация контроллеров
	pvzCtrl := controllers.NewPVZController(createPVZUC, listPVZsUC, closeReceptionUC, deleteLastProductUC)
	receptionCtrl := controllers.NewReceptionController(createReceptionUC)
	productCtrl := controllers.NewProductController(addProductUC)
	authCtrl := controllers.NewAuthController(dummyLoginUC, nil, nil)

	// Настройка роутера
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Публичные роуты
	r.POST("/dummyLogin", authCtrl.DummyLogin)

	// Защищенные роуты
	auth := r.Group("/", controllers.JWTAuthMiddleware("test_secret"))
	auth.POST("/pvz", pvzCtrl.Create)
	auth.POST("/receptions", receptionCtrl.Create)
	auth.POST("/products", productCtrl.Add)
	auth.POST("/pvz/:pvzId/close_last_reception", pvzCtrl.CloseLastReception)
	auth.POST("/pvz/:pvzId/delete_last_product", pvzCtrl.DeleteLastProduct) // Добавляем эндпоинт для удаления товара

	return r, db
}

func TestPVZAPIIntegration(t *testing.T) {
	// Arrange: настраиваем тестовый сервер и базу данных
	r, db := setupTestServer(t)
	ctx := context.Background()

	// Act: выполняем сценарий тестирования

	// 1. Получаем токен модератора
	moderatorToken := getToken(t, r, entities.UserRoleModerator)

	// 2. Создаем ПВЗ
	pvzID := createPVZ(t, r, moderatorToken)

	// 3. Получаем токен сотрудника ПВЗ
	staffToken := getToken(t, r, entities.UserRolePVZStaff)

	// 4. Создаем приёмку
	receptionID := createReception(t, r, staffToken, pvzID)

	// 5. Добавляем 50 товаров
	for i := 0; i < 50; i++ {
		productType := entities.ProductElectronics
		if i%3 == 0 {
			productType = entities.ProductClothes
		} else if i%3 == 1 {
			productType = entities.ProductShoes
		}
		addProduct(t, r, staffToken, pvzID, productType)
	}

	// 6. Закрываем приёмку
	closeReception(t, r, staffToken, pvzID)

	// Assert: проверяем результаты
	productRepo := repositories.NewPGProductRepository(db)
	products, err := productRepo.ListByReception(ctx, receptionID)
	require.NoError(t, err)
	require.Len(t, products, 50)
}

func getToken(t *testing.T, r *gin.Engine, role entities.UserRole) string {
	body := map[string]string{"role": string(role)}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Contains(t, response, "token")
	return response["token"]
}

func createPVZ(t *testing.T, r *gin.Engine, token string) uuid.UUID {
	body := map[string]string{"city": string(entities.CityMoscow)}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var response entities.PVZ
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	return response.ID
}

func createReception(t *testing.T, r *gin.Engine, token string, pvzID uuid.UUID) uuid.UUID {
	body := map[string]string{"pvzId": pvzID.String()}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response entities.Reception
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	return response.ID
}

func addProduct(t *testing.T, r *gin.Engine, token string, pvzID uuid.UUID, productType entities.ProductType) {
	body := map[string]interface{}{
		"pvzId": pvzID.String(),
		"type":  string(productType),
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func closeReception(t *testing.T, r *gin.Engine, token string, pvzID uuid.UUID) {
	req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/close_last_reception", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
