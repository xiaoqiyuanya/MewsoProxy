package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"mewsoproxy/server/pkg/apperror"
	redisc "mewsoproxy/server/pkg/redis"
	"mewsoproxy/server/pkg/response"
)

func RateLimit(rds *redisc.Client, bucket string, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := redisc.RateLimitKey(bucket, c.ClientIP())
		ctx := c.Request.Context()
		val, err := rds.Incr(ctx, key)
		if err != nil {
			response.Fail(c, apperror.CodeRateLimited, "请求过于频繁")
			c.Abort()
			return
		}
		if val == 1 {
			_ = rds.Expire(ctx, key, window)
		}
		if val > int64(limit) {
			response.Fail(c, apperror.CodeRateLimited, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}

func IncrKey(ctx context.Context, rds *redisc.Client, key string, window time.Duration) (int64, error) {
	val, err := rds.Incr(ctx, key)
	if err != nil {
		return 0, err
	}
	if val == 1 {
		_ = rds.Expire(ctx, key, window)
	}
	return val, nil
}
