package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ydoro/wishlist/internal/domain"
)

func SetupRoutes(
	r *gin.Engine,
	customerCreation domain.CreateCustomerUC,
	userAuthentication domain.Authenticator,
	authMiddleware gin.HandlerFunc,
	customerGetter domain.ShowCustomerDataUC,
	customerUpdater domain.UpdateCustomerUC,
	customerDeleter domain.DeleteCustomerUC,
	wishlistCreator domain.CreateWishlistUseCase,
	wishlistDeleter domain.DeleteWishlistUseCase,
	wishlistGetter domain.ShowWishlistUseCase,
	wishlistUpdater domain.UpdateWishListUseCase,
	productGetter domain.GetProductUseCase,
	productLister domain.ListProductsUseCase,

) *gin.Engine {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	NewAuthHandler(api, userAuthentication)
	NewProductHandler(api, productGetter, productLister)

	customerRoutes := NewCustomerHandler(api, customerCreation, authMiddleware, customerGetter, customerUpdater, customerDeleter)
	SetupWishlistHandler(
		customerRoutes,
		authMiddleware,
		wishlistCreator,
		wishlistDeleter,
		wishlistGetter,
		wishlistUpdater,
	)

	return r
}
