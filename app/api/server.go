package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/muasx88/stealth_commerce/app/config/middlewares"
	"github.com/muasx88/stealth_commerce/app/modules/auth"
	"github.com/muasx88/stealth_commerce/app/modules/cart"
	"github.com/muasx88/stealth_commerce/app/modules/checkout"
	"github.com/muasx88/stealth_commerce/app/modules/product"

	"github.com/gin-contrib/requestid"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	router *gin.Engine
	dbConn *pgxpool.Pool
}

func NewServer(dbConn *pgxpool.Pool) *Server {
	server := &Server{
		dbConn: dbConn,
	}

	server.setupRoute()
	return server
}

func (server *Server) setupRoute() {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})

	router.Use(requestid.New())
	router.Use(gin.Recovery())
	router.Use(middlewares.LoggingMiddleware())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "stealth-commerce api",
			"version": "1.0.0",
		})
	})

	auth.AuthRoute(router, server.dbConn)
	product.ProductRoute(router, server.dbConn)
	cart.CartRoute(router, server.dbConn)
	checkout.CheckoutRoute(router, server.dbConn)

	server.router = router
}

func (server *Server) Start(ctx context.Context, address string) {
	srv := &http.Server{
		Addr:    ":" + address,
		Handler: server.router,
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- srv.ListenAndServe()

	}()

	shutDownChann := make(chan os.Signal, 1)
	signal.Notify(shutDownChann, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-shutDownChann:
		log.Println("shutdown signal", sig)
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal(err)
			srv.Close()
		}
	case err := <-serverErr:
		if err != nil {
			log.Fatalf("server : %v", err)
		}
	}

	log.Printf("http server run on port %s", address)
}
