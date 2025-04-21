package main

import (
	"context"
	"database/sql"
	"log"
	"os"
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
)

// --- Адаптеры для DI ---
type userRepoForRegister struct{ *repositories.PGUserRepository }

func (a *userRepoForRegister) Create(ctx context.Context, user entities.User, passwordHash string) (entities.User, error) {
	return a.PGUserRepository.Create(ctx, user, passwordHash)
}
func (a *userRepoForRegister) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	u, _, err := a.PGUserRepository.GetByEmail(ctx, email)
	return u, err
}

type productRepoForDelete struct {
	*repositories.PGProductRepository
}

func (a *productRepoForDelete) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := a.PGProductRepository.DeleteLast(ctx, id)
	return err
}

type receptionRepoForClose struct {
	*repositories.PGReceptionRepository
}

func (a *receptionRepoForClose) GetActive(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
	return a.PGReceptionRepository.GetActive(ctx, pvzID)
}
func (a *receptionRepoForClose) Save(ctx context.Context, rec entities.Reception) (entities.Reception, error) {
	return a.PGReceptionRepository.Save(ctx, rec)
}
func (a *receptionRepoForClose) CloseLast(ctx context.Context, pvzID uuid.UUID) (entities.Reception, error) {
	err := a.PGReceptionRepository.CloseLast(ctx, pvzID, time.Now().UTC())
	if err != nil {
		return entities.Reception{}, err
	}
	return entities.Reception{}, nil // TODO: вернуть реальный rec, если нужно
}

type productRepoForList struct {
	*repositories.PGProductRepository
}

func (a *productRepoForList) Save(ctx context.Context, p entities.Product) (entities.Product, error) {
	return a.PGProductRepository.Save(ctx, p)
}
func (a *productRepoForList) DeleteLast(ctx context.Context, id uuid.UUID) error {
	_, err := a.PGProductRepository.DeleteLast(ctx, id)
	return err
}
func (a *productRepoForList) ListByReception(ctx context.Context, recID uuid.UUID) ([]entities.Product, error) {
	return a.PGProductRepository.ListByReception(ctx, recID)
}

type receptionRepoForList struct {
	*repositories.PGReceptionRepository
}

func (a *receptionRepoForList) GetActive(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
	return a.PGReceptionRepository.GetActive(ctx, pvzID)
}
func (a *receptionRepoForList) Save(ctx context.Context, rec entities.Reception) (entities.Reception, error) {
	return a.PGReceptionRepository.Save(ctx, rec)
}
func (a *receptionRepoForList) CloseLast(ctx context.Context, pvzID uuid.UUID) (entities.Reception, error) {
	err := a.PGReceptionRepository.CloseLast(ctx, pvzID, time.Now().UTC())
	if err != nil {
		return entities.Reception{}, err
	}
	return entities.Reception{}, nil // TODO: вернуть реальный rec, если нужно
}

func main() {
	_ = godotenv.Load()
	cfg := configs.LoadConfig()

	// --- Postgres ---
	db, err := sql.Open("pgx", cfg.PGDSN)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	// --- Репозитории ---
	userRepo := repositories.NewPGUserRepository(db)
	pvzRepo := repositories.NewPGPVZRepository(db)
	receptionRepo := repositories.NewPGReceptionRepository(db)
	productRepo := repositories.NewPGProductRepository(db)

	// --- Usecase ---
	dummyLoginUC := usecases.NewDummyLoginUseCase(cfg)
	registerUC := usecases.NewRegisterUseCase(&userRepoForRegister{userRepo})
	loginUC := usecases.NewLoginUseCase(userRepo, cfg)
	createPVZUC := usecases.NewCreatePVZUseCase(pvzRepo)
	listPVZsUC := usecases.NewListPVZsUseCase(pvzRepo, &receptionRepoForList{receptionRepo}, &productRepoForList{productRepo})
	closeReceptionUC := usecases.NewCloseReceptionUseCase(&receptionRepoForClose{receptionRepo})
	deleteLastProductUC := usecases.NewDeleteLastProductUseCase(&productRepoForDelete{productRepo}, &receptionRepoForClose{receptionRepo})
	addProductUC := usecases.NewAddProductUseCase(productRepo, receptionRepo)
	createReceptionUC := usecases.NewCreateReceptionUseCase(receptionRepo)

	// --- Контроллеры ---
	authCtrl := controllers.NewAuthController(dummyLoginUC, registerUC, loginUC)
	pvzCtrl := controllers.NewPVZController(createPVZUC, listPVZsUC, closeReceptionUC, deleteLastProductUC)
	productCtrl := controllers.NewProductController(addProductUC)
	receptionCtrl := controllers.NewReceptionController(createReceptionUC)

	r := gin.Default()

	// --- Auth ---
	r.POST("/dummyLogin", authCtrl.DummyLogin)
	r.POST("/register", authCtrl.Register)
	r.POST("/login", authCtrl.Login)

	// --- Защищённые эндпоинты ---
	authMW := controllers.JWTAuthMiddleware(cfg.JWTSecret)
	pvz := r.Group("/pvz", authMW)
	pvz.POST("/", pvzCtrl.Create)
	pvz.GET("/", pvzCtrl.List)
	pvz.POST("/:pvzId/close_last_reception", pvzCtrl.CloseLastReception)
	pvz.POST("/:pvzId/delete_last_product", pvzCtrl.DeleteLastProduct)

	product := r.Group("/products", authMW)
	product.POST("/", productCtrl.Add)

	reception := r.Group("/receptions", authMW)
	reception.POST("/", receptionCtrl.Create)

	// Healthcheck
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting HTTP server on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
