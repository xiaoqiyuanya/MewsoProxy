package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"

	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/response"
	"mewsoproxy/server/pkg/token"
)

type ctxKey string

const (
	CtxUserID ctxKey = "user_id"
	CtxUserRole ctxKey = "user_role"
	CtxJTI ctxKey = "jti"
)

type TokenValidator interface {
	IsBlacklisted(ctx context.Context, jti string) (bool, error)
}

func Auth(secret string, validator TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			if q := c.Query("token"); q != "" {
				auth = "Bearer " + q
			}
		}
		if auth == "" {
			response.Fail(c, apperror.CodeNotLogin, "未登录")
			c.Abort()
			return
		}
		if !strings.HasPrefix(auth, "Bearer ") {
			response.Fail(c, apperror.CodeNotLogin, "登录凭证无效")
			c.Abort()
			return
		}
		raw := strings.TrimPrefix(auth, "Bearer ")
		claims, err := token.ParseAccessToken(raw, secret)
		if err != nil {
			response.Fail(c, apperror.CodeTokenExpired, "登录已过期，请重新登录")
			c.Abort()
			return
		}
		if validator != nil {
			blk, _ := validator.IsBlacklisted(c.Request.Context(), claims.JTI)
			if blk {
				response.Fail(c, apperror.CodeTokenExpired, "登录已失效")
				c.Abort()
				return
			}
		}
		c.Set(string(CtxUserID), claims.UserID)
		c.Set(string(CtxUserRole), claims.Role)
		c.Set(string(CtxJTI), claims.JTI)
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(string(CtxUserRole))
		if role != "admin" {
			response.Fail(c, apperror.CodeNoPermission, "无权限操作")
			c.Abort()
			return
		}
		c.Next()
	}
}

func UserIDFrom(c *gin.Context) uint {
	v, ok := c.Get(string(CtxUserID))
	if !ok {
		return 0
	}
	id, _ := v.(uint)
	return id
}

func JTIFrom(c *gin.Context) string {
	v, _ := c.Get(string(CtxJTI))
	s, _ := v.(string)
	return s
}
