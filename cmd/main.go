package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ydoro/wishlist/config"
	"github.com/ydoro/wishlist/internal/infra/adapter"
	postgresDB "github.com/ydoro/wishlist/internal/infra/db/postgres"
	"github.com/ydoro/wishlist/internal/infra/delivery/http"
	"github.com/ydoro/wishlist/internal/infra/delivery/http/middleware"
	"github.com/ydoro/wishlist/internal/usecase"

	_ "github.com/ydoro/wishlist/docs" // for swagger
)

// @title Wishlist API GO
// @version 1.0
// @description A powerful API for managing customers wishlists.
// @host localhost:8080
// @BasePath /api/
// @securityDefinitions.apikey BearerAuth
// @in Header
// @name Authorization
func main() {
	cfg := config.LoadConfig()
	r := gin.Default()

	// here we can use some DI framework or some factory to create the use cases
	conn, err := postgresDB.Connect(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBSSL,
	)

	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}

	defer conn.Close()

	customerRepo := postgresDB.NewCustomerRepository(conn)
	idGenerator := adapter.UUIDGenerator{}
	hasher := adapter.NewPasswordHasher(10)
	jwtEcnoder := adapter.NewJWTEncrypter(cfg.JWTSecret)

	authMiddleware := middleware.NewAuthMiddleware(jwtEcnoder)

	ucs := usecase.NewCreateCustomerUseCase(customerRepo, idGenerator, adapter.NewPasswordHasher(10))
	showCustomerUC := usecase.NewGetCustomerData(customerRepo)
	authUC := usecase.NewPasswordAuthenticationUseCase(hasher, customerRepo, jwtEcnoder)
	updateCustomerUc := usecase.NewUpdateCustomerUseCase(customerRepo, customerRepo, customerRepo)

	router := http.SetupRoutes(r, ucs, authUC, authMiddleware.Handle, showCustomerUC, updateCustomerUc)
	router.Run(":8080")
}
