package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	config "github.com/ydoro/wishlist/config/customer"
	"github.com/ydoro/wishlist/internal/customer/infra/adapter"
	postgresDB "github.com/ydoro/wishlist/internal/customer/infra/db/postgres"
	"github.com/ydoro/wishlist/internal/customer/infra/delivery/http"
	"github.com/ydoro/wishlist/internal/customer/usecase"
)

// @title Wishlist API GO
// @version 1.0
// @description A powerful API for managing customers wishlists.
// @host localhost:8080
// @BasePath /api/v1
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

	ucs := usecase.NewCreateCustomerUseCase(customerRepo, idGenerator)

	router := http.SetupRoutes(r, ucs, http.AuthMiddleware(cfg.JWTSecret))
	router.Run(":8080")
}
