package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware creates a structured logging middleware
func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate correlation ID
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.New().String()
		}
		c.Set("correlation_id", correlationID)
		c.Header("X-Correlation-ID", correlationID)

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		logEntry := logger.WithFields(logrus.Fields{
			"correlation_id": correlationID,
			"method":         c.Request.Method,
			"path":           c.Request.URL.Path,
			"query":          c.Request.URL.RawQuery,
			"status":         statusCode,
			"duration_ms":    duration.Milliseconds(),
			"client_ip":      c.ClientIP(),
			"user_agent":     c.Request.UserAgent(),
		})

		if username, exists := c.Get("username"); exists {
			logEntry = logEntry.WithField("username", username)
		}

		if statusCode >= 500 {
			logEntry.Error("Request failed")
		} else if statusCode >= 400 {
			logEntry.Warn("Request returned client error")
		} else {
			logEntry.Info("Request completed")
		}
	}
}

// RecoveryMiddleware recovers from panics and logs them
func RecoveryMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				correlationID, _ := c.Get("correlation_id")
				logger.WithFields(logrus.Fields{
					"correlation_id": correlationID,
					"error":          err,
					"path":           c.Request.URL.Path,
				}).Error("Panic recovered")

				c.JSON(500, gin.H{
					"error":          "internal server error",
					"correlation_id": correlationID,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
