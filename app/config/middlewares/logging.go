package middlewares

import (
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Starting time
		startTime := time.Now()

		// Processing request
		c.Next()

		// End Time
		endTime := time.Now()

		// execution time
		latencyTime := endTime.Sub(startTime)

		// Request method
		reqMethod := c.Request.Method

		// Request route
		reqUri := c.Request.RequestURI

		// status code
		statusCode := c.Writer.Status()

		// Request IP
		clientIP := c.ClientIP()

		requestId := requestid.Get(c)

		log.WithFields(log.Fields{
			"METHOD":     reqMethod,
			"URI":        reqUri,
			"STATUS":     statusCode,
			"LATENCY":    latencyTime,
			"CLIENT_IP":  clientIP,
			"REQUEST_ID": requestId,
		}).Info("HTTP REQUEST")

		c.Next()
	}
}
