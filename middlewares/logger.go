// middlewares/logger.go
package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process Request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get status
		status := c.Writer.Status()

		// Log format
		logger.WithFields(logrus.Fields{
			"status":        status,
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"ip":            c.ClientIP(),
			"latency":       latency,
			"user_agent":    c.Request.UserAgent(),
			"error_message": c.Errors.ByType(gin.ErrorTypePrivate).String(),
		}).Info("Request completed")
	}
}
