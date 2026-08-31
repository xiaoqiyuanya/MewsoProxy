package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"mewsoproxy/server/pkg/apperror"
	redisc "mewsoproxy/server/pkg/redis"
	"mewsoproxy/server/pkg/response"
)

func Idempotent(rds *redisc.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")
		if requestID != "" {
			key := redisc.IdempotentKey(requestID)
			ok, err := rds.SetNX(c.Request.Context(), key, "1", 24*time.Hour)
			if err == nil && !ok {
				response.Fail(c, apperror.CodeParamFormat, "请勿重复提交请求")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
