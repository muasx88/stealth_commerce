package dbhandler

import (
	"time"

	"github.com/muasx88/stealth_commerce/app/internals/config"
)

func getMaxConn() int32 {
	defaultMaxConns := 3
	maxConnEnv := config.Config.DBMaxConnections

	if maxConnEnv > 0 {
		defaultMaxConns = maxConnEnv
	}

	return int32(defaultMaxConns)
}

func getMaxIdleTime() time.Duration {
	defaultIdleTime := 120
	idleTimeEnv := config.Config.DBMaxIdleConnections

	if idleTimeEnv > 0 {
		defaultIdleTime = idleTimeEnv
	}

	return time.Duration(defaultIdleTime)
}
