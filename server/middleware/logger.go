package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"log/slog"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		slog.Info("http request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.String("ip", c.ClientIP()),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", time.Since(start)),
		)
	}
}
