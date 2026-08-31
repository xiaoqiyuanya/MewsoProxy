package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"mewsoproxy/server/pkg/apperror"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered", "err", err, "path", c.FullPath())
				c.JSON(http.StatusOK, gin.H{
					"code":    apperror.CodeDBError,
					"message": "系统错误",
					"data":    nil,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
