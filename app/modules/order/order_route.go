package order

import (
	"github.com/muasx88/stealth_commerce/app/config/middlewares"
	handler "github.com/muasx88/stealth_commerce/app/modules/order/delivery/http"
	"github.com/muasx88/stealth_commerce/app/modules/order/repository"
	"github.com/muasx88/stealth_commerce/app/modules/order/usecase"

	buyerRepo "github.com/muasx88/stealth_commerce/app/modules/buyer/repository"
	buyerUsecase "github.com/muasx88/stealth_commerce/app/modules/buyer/usecase"

	cartRepo "github.com/muasx88/stealth_commerce/app/modules/cart/repository"
	cartUsecase "github.com/muasx88/stealth_commerce/app/modules/cart/usecase"
	productRepo "github.com/muasx88/stealth_commerce/app/modules/product/repository"
	productUsecase "github.com/muasx88/stealth_commerce/app/modules/product/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func OrderRoute(r *gin.Engine, db *pgxpool.Pool) {
	buyerRepo := buyerRepo.NewBuyerRepository(db)
	cartRepo := cartRepo.NewCartRepository(db)
	productRepo := productRepo.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	buyerUsecase := buyerUsecase.NewProductUsecase(buyerRepo)
	cartUsecase := cartUsecase.NewCartUsecase(cartRepo)
	productUsecase := productUsecase.NewProductUsecase(productRepo)
	orderUsecase := usecase.NewOrderUseOrderUsecase(orderRepo, cartUsecase, productUsecase)

	handler := handler.NewOrderHandler(orderUsecase)

	orderRoute := r.Group("/orders").Use(middlewares.BasicAuth(buyerUsecase))
	orderRoute.POST("", handler.CreateOrder)

}
