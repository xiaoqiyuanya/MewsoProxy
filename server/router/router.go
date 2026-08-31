package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mewsoproxy/server/config"
	"mewsoproxy/server/handler/admin"
	"mewsoproxy/server/handler/auth"
	"mewsoproxy/server/handler/node"
	"mewsoproxy/server/handler/order"
	"mewsoproxy/server/handler/payment"
	"mewsoproxy/server/handler/plan"
	"mewsoproxy/server/handler/subscribe"
	"mewsoproxy/server/handler/user"
	"mewsoproxy/server/middleware"
	redisc "mewsoproxy/server/pkg/redis"
	"mewsoproxy/server/pkg/response"
	adminsvc "mewsoproxy/server/service/admin"
	authsvc "mewsoproxy/server/service/auth"
	ordersvc "mewsoproxy/server/service/order"
	nodesvc "mewsoproxy/server/service/node"
	paymentsvc "mewsoproxy/server/service/payment"
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
	paymentSvc := paymentsvc.NewService(db, cfg, orderSvc)

	authHandler := auth.New(authSvc, userSvc)
	userHandler := user.New(userSvc)
	planHandler := plan.New(planSvc)
	orderHandler := order.New(orderSvc)
	paymentHandler := payment.New(paymentSvc)
	subSvc := subsvc.New(cfg, db)
	subHandler := subscribe.New(subSvc, userSvc)

	adminSvc := adminsvc.New(db, cfg, rds, orderSvc)
	adminHandler := admin.New(adminSvc, authSvc, userSvc)

	nodeSvc := nodesvc.New(db, cfg)
	nodeHandler := node.New(nodeSvc)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS([]string{"http://localhost:5173", "http://localhost:3000"}))

	r.GET("/healthz", func(c *gin.Context) {
		response.OK(c, gin.H{"status": "ok"})
	})
	r.GET("/subscribe", subHandler.Download)

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

		authed.GET("/payment/channels", paymentHandler.Channels)
		authed.POST("/payment/create", paymentHandler.Create)
	}

	api.POST("/payment/notify", paymentHandler.Notify)
	api.GET("/payment/notify", paymentHandler.Notify)

	api.POST("/node/report", nodeHandler.Report)

	api.POST("/admin/login", adminHandler.Login)

	adminGroup := api.Group("/admin", middleware.Auth(cfg.JWT.AccessSecret, authSvc), middleware.AdminOnly())
	{
		adminGroup.GET("/config/fetch", adminHandler.GetConfig)
		adminGroup.POST("/config/save", adminHandler.SaveConfig)
		adminGroup.GET("/system/status", adminHandler.SystemStatus)

		adminGroup.GET("/plan/list", adminHandler.ListPlans)
		adminGroup.POST("/plan/save", adminHandler.SavePlan)
		adminGroup.POST("/plan/drop", adminHandler.DropPlan)

		adminGroup.GET("/user/list", adminHandler.ListUsers)
		adminGroup.POST("/user/info", adminHandler.GetUser)
		adminGroup.POST("/user/update", adminHandler.UpdateUser)
		adminGroup.POST("/user/ban", adminHandler.BanUser)
		adminGroup.POST("/user/reset_secret", adminHandler.ResetSecret)

		adminGroup.GET("/order/list", adminHandler.ListOrders)
		adminGroup.GET("/order/info", adminHandler.GetOrder)
		adminGroup.POST("/order/cancel", adminHandler.CancelOrder)
		adminGroup.POST("/order/paid", adminHandler.MarkOrderPaid)

		adminGroup.GET("/server/group/list", adminHandler.ListGroups)
		adminGroup.POST("/server/group/save", adminHandler.SaveGroup)
		adminGroup.POST("/server/group/drop", adminHandler.DropGroup)
		adminGroup.GET("/server/node/list", adminHandler.ListNodes)
		adminGroup.POST("/server/node/save", adminHandler.SaveNode)
		adminGroup.POST("/server/node/drop", adminHandler.DropNode)
		adminGroup.POST("/server/node/install", adminHandler.NodeInstall)
		adminGroup.GET("/server/node/install/log", adminHandler.NodeInstallLog)

		adminGroup.GET("/coupon/list", adminHandler.ListCoupons)
		adminGroup.POST("/coupon/save", adminHandler.SaveCoupon)
		adminGroup.POST("/coupon/drop", adminHandler.DropCoupon)
		adminGroup.POST("/coupon/show", adminHandler.ToggleCouponShow)

		adminGroup.GET("/notice/list", adminHandler.ListNotices)
		adminGroup.POST("/notice/save", adminHandler.SaveNotice)
		adminGroup.POST("/notice/drop", adminHandler.DropNotice)
		adminGroup.POST("/notice/show", adminHandler.ToggleNoticeShow)

		adminGroup.GET("/payment/list", adminHandler.ListPayments)
		adminGroup.POST("/payment/save", adminHandler.SavePayment)
		adminGroup.POST("/payment/drop", adminHandler.DropPayment)
		adminGroup.POST("/payment/show", adminHandler.TogglePaymentShow)
	}

	return r
}
