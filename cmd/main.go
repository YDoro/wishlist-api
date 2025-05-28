package main

import (
	"fmt"
	https "net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydoro/wishlist/config"
	"github.com/ydoro/wishlist/internal/infra/adapter"
	postgresDB "github.com/ydoro/wishlist/internal/infra/db/postgres"
	"github.com/ydoro/wishlist/internal/infra/delivery/http"
	"github.com/ydoro/wishlist/internal/infra/delivery/http/middleware"
	"github.com/ydoro/wishlist/internal/infra/services"
	"github.com/ydoro/wishlist/internal/usecase"

	_ "github.com/ydoro/wishlist/docs" // for swagger
)

// @title Wishlist API GO
// @version 1.0
// @description A powerful API for managing customers wishlists.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in Header
// @name Authorization
func main() {
	cfg := config.LoadConfig()
	r := gin.Default()

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

	// Services
	redisDb, err := strconv.Atoi(cfg.CACHE_DATABASE)
	if err != nil {
		redisDb = 0
	}

	rediscfg := &adapter.RedisCacheConfig{
		Addr:     cfg.CACHE_URL,
		Password: cfg.CACHE_PASSWORD,
		DB:       redisDb,
	}
	redis := adapter.NewRedisCache(*rediscfg)
	httpClient := &https.Client{}

	productService := services.NewFakeProductAPIService(cfg.PRODUCT_API_URL, httpClient)

	// TODO - improve DI, use a factory or a DI framework
	customerRepo := postgresDB.NewCustomerRepository(conn)
	wishlistRepo := postgresDB.NewWishlistRepository(conn)
	productRepo := postgresDB.NewProductRepository(conn)
	idGenerator := adapter.UUIDGenerator{}
	hasher := adapter.NewPasswordHasher(10)
	jwtEcnoder := adapter.NewJWTEncrypter(cfg.JWTSecret)

	authMiddleware := middleware.NewAuthMiddleware(jwtEcnoder)

	authUC := usecase.NewPasswordAuthenticationUseCase(hasher, customerRepo, jwtEcnoder)

	createCustomerUC := usecase.NewCreateCustomerUseCase(customerRepo, idGenerator, adapter.NewPasswordHasher(10))
	showCustomerUC := usecase.NewGetCustomerData(customerRepo)
	updateCustomerUc := usecase.NewUpdateCustomerUseCase(customerRepo, customerRepo, customerRepo)
	deleteCustomerUc := usecase.NewDeleteCustomerUseCase(customerRepo, customerRepo)

	getProductUc := usecase.NewGetProductAndStoreIfNeededUseCase(cfg.CACHE_TTL, redis, productService, productRepo, productRepo, productRepo)
	listProductUc := usecase.NewListProductsAndStoreUseCase(cfg.CACHE_TTL, redis, productService, productRepo, productRepo, productRepo)

	createWishlistUc := usecase.NewCreateWishlistUseCase(wishlistRepo, wishlistRepo, customerRepo, idGenerator)
	deleteWishlistUc := usecase.NewDeleteWishlistUseCase(customerRepo, wishlistRepo, wishlistRepo)
	getWishlistUC := usecase.NewShowWishlistUseCase(wishlistRepo, customerRepo, getProductUc)
	updateWishlistUC := usecase.NewUpdateWishListUseCase(customerRepo, wishlistRepo, wishlistRepo, getProductUc)
	listWishlistUC := usecase.NewListCustomerWishlistsUseCase(customerRepo, wishlistRepo, getProductUc)

	router := http.SetupRoutes(
		r,
		createCustomerUC,
		authUC,
		authMiddleware.Handle,
		showCustomerUC,
		updateCustomerUc,
		deleteCustomerUc,
		createWishlistUc,
		deleteWishlistUc,
		getWishlistUC,
		updateWishlistUC,
		getProductUc,
		listProductUc,
		listWishlistUC,
	)

	router.Run(":8080")
}
