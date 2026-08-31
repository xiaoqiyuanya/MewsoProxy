package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mewsoproxy/server/config"
	"mewsoproxy/server/handler/auth"
	"mewsoproxy/server/handler/order"
	"mewsoproxy/server/handler/payment"
	"mewsoproxy/server/handler/plan"
	"mewsoproxy/server/handler/subscribe"
	"mewsoproxy/server/handler/user"
	"mewsoproxy/server/middleware"
	redisc "mewsoproxy/server/pkg/redis"
	"mewsoproxy/server/pkg/response"
	authsvc "mewsoproxy/server/service/auth"
	ordersvc "mewsoproxy/server/service/order"
	plansvc "mewsoproxy/server/service/plan"
	subsvc "mewsoproxy/server/service/subscribe"
	usersvc "mewsoproxy/server/service/user"
)

func New(cfg *config.Config, db *gorm.DB, rds *redisc.Client) *gin.Engine {
	if gin.Mode() == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	}

	userSvc := usersvc.New(db, cfg)
	authSvc := authsvc.New(cfg, rds, userSvc)
	planSvc := plansvc.New(db)
	orderSvc := ordersvc.New(db, cfg)

	authHandler := auth.New(authSvc, userSvc)
	userHandler := user.New(userSvc)
	planHandler := plan.New(planSvc)
	orderHandler := order.New(orderSvc)
	paymentHandler := payment.New(orderSvc)
	subSvc := subsvc.New(cfg)
	subHandler := subscribe.New(subSvc, userSvc)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS([]string{"http://localhost:5173", "http://localhost:3000"}))

	r.GET("/healthz", func(c *gin.Context) {
		response.OK(c, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")

	api.POST("/auth/register", middleware.RateLimit(rds, "register", 5, time.Minute), authHandler.Register)
	api.POST("/auth/login", middleware.RateLimit(rds, "login", 5, time.Minute), authHandler.Login)
	api.POST("/auth/refresh", authHandler.Refresh)
	api.GET("/plan/list", planHandler.List)

	authed := api.Group("", middleware.Auth(cfg.JWT.AccessSecret, authSvc))
	{
		authed.GET("/user/me", userHandler.Me)
		authed.POST("/user/logout", userHandler.Logout)
		authed.GET("/user/subscribe", subHandler.Info)

		authed.POST("/order/create", middleware.Idempotent(rds), orderHandler.Create)
		authed.GET("/order/list", orderHandler.List)
		authed.POST("/order/detail", orderHandler.Detail)
		authed.POST("/order/cancel", orderHandler.Cancel)
	}

	api.POST("/payment/notify", paymentHandler.Notify)

	return r
}
