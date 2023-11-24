package main

import (
	"context"

	"github.com/muasx88/stealth_commerce/app/api"
	"github.com/muasx88/stealth_commerce/app/internals/config"
	"github.com/muasx88/stealth_commerce/app/internals/dbhandler"

	log "github.com/sirupsen/logrus"
)

func init() {
	err := config.LoadConfig("./")
	if err != nil {
		log.Fatalf("error load config %v", err)
	}

	// setup logrus
	logLevel, err := log.ParseLevel(config.Config.LOG_LEVEL)
	if err != nil {
		logLevel = log.InfoLevel
	}

	log.SetLevel(logLevel)
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	ctx := context.Background()
	dbConn, err := dbhandler.GetDBConnection(ctx)
	if err != nil {
		log.Fatalf("error connect db: %v", err)
	}

	srv := api.NewServer(dbConn)
	srv.Start(ctx, config.Config.PORT)
}
