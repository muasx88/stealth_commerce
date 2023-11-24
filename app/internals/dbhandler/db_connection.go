package dbhandler

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/muasx88/stealth_commerce/app/internals/config"
)

var DBConn *pgxpool.Pool

func createDbConnection(ctx context.Context) (*pgxpool.Pool, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable binary_parameters=yes TimeZone=Asia/Jakarta",
		config.Config.DBHost,
		config.Config.DBPort,
		config.Config.DBUsername,
		config.Config.DBPassword,
		config.Config.DBName,
	)

	// Create the connection pool
	poolConfig, err := pgxpool.ParseConfig(psqlconn)
	if err != nil {
		return nil, err
	}

	// Customize the pool configuration if needed
	poolConfig.MaxConns = getMaxConn()
	poolConfig.MaxConnIdleTime = getMaxIdleTime() * time.Second
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeDescribeExec
	poolConfig.ConnConfig.RuntimeParams = map[string]string{}
	poolConfig.ConnConfig.RuntimeParams["application_name"] = "jb_chat_api"

	dbPool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	err = dbPool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return dbPool, nil
}

func GetDBConnection(ctx context.Context) (*pgxpool.Pool, error) {
	if DBConn == nil {
		var err error
		DBConn, err = createDbConnection(ctx)
		if err != nil {
			return nil, err
		}
	}

	return DBConn, nil
}
