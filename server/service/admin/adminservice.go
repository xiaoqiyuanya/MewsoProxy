package admin

import (
	"gorm.io/gorm"

	"mewsoproxy/server/config"
	"mewsoproxy/server/pkg/redis"
	ordersvc "mewsoproxy/server/service/order"
)

type Service struct {
	db  *gorm.DB
	cfg *config.Config
	rds *redis.Client
	orderSvc *ordersvc.Service
}

func New(db *gorm.DB, cfg *config.Config, rds *redis.Client, orderSvc *ordersvc.Service) *Service {
	return &Service{db: db, cfg: cfg, rds: rds, orderSvc: orderSvc}
}
