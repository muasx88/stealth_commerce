package cart

import (
	"github.com/muasx88/stealth_commerce/app/config/middlewares"
	buyerRepo "github.com/muasx88/stealth_commerce/app/modules/buyer/repository"
	buyerUsecase "github.com/muasx88/stealth_commerce/app/modules/buyer/usecase"
	handler "github.com/muasx88/stealth_commerce/app/modules/cart/delivery/http"
	"github.com/muasx88/stealth_commerce/app/modules/cart/repository"
	"github.com/muasx88/stealth_commerce/app/modules/cart/usecase"
	productRepo "github.com/muasx88/stealth_commerce/app/modules/product/repository"
	prouctUsecase "github.com/muasx88/stealth_commerce/app/modules/product/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CartRoute(r *gin.Engine, db *pgxpool.Pool) {
	buyerRepo := buyerRepo.NewBuyerRepository(db)
	cartRepo := repository.NewCartRepository(db)
	productRepo := productRepo.NewProductRepository(db)

	buyerUsecase := buyerUsecase.NewProductUsecase(buyerRepo)
	cartUsecase := usecase.NewCartUsecase(cartRepo)
	prouctUsecase := prouctUsecase.NewProductUsecase(productRepo)

	handler := handler.NewCartHandler(cartUsecase, prouctUsecase)

	cartRoute := r.Group("/carts").Use(middlewares.BasicAuth(buyerUsecase))
	cartRoute.GET("", handler.GetAll)
	cartRoute.POST("", handler.Add)
	cartRoute.PUT("/:id", handler.Update)
	cartRoute.DELETE("/:id", handler.Delete)
}
