package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ydoro/wishlist/internal/domain"
)

func SetupRoutes(r *gin.Engine, uc domain.CreateCustomerUC) *gin.Engine {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	NewCustomerHandler(api, uc)

	return r
}
