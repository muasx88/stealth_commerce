package product

import (
	"github.com/muasx88/stealth_commerce/app/config/middlewares"
	handler "github.com/muasx88/stealth_commerce/app/modules/product/delivery/http"
	"github.com/muasx88/stealth_commerce/app/modules/product/repository"
	"github.com/muasx88/stealth_commerce/app/modules/product/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ProductRoute(r *gin.Engine, db *pgxpool.Pool) {
	productRepo := repository.NewProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUsecase)

	r.GET("/products", productHandler.GetAll)
	r.GET("/products/:id", productHandler.Detail)

	productRoute := r.Group("/products").Use(middlewares.JwtProtected())
	productRoute.POST("", productHandler.Add)
	productRoute.PUT("/:id", productHandler.Update)
	productRoute.DELETE("/:id", productHandler.Delete)
}
