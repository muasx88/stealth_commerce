package auth

import (
	handler "github.com/muasx88/stealth_commerce/app/modules/auth/delivery/http"
	"github.com/muasx88/stealth_commerce/app/modules/auth/repository"
	"github.com/muasx88/stealth_commerce/app/modules/auth/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AuthRoute(r *gin.Engine, db *pgxpool.Pool) {
	repo := repository.NewAuthRepository(db)
	usecase := usecase.NewAuthUsecase(repo)
	handler := handler.NewAuthHandler(usecase)

	r.POST("/login", handler.Login)
}
