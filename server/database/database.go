package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"mewsoproxy/server/config"
	"mewsoproxy/server/model"
)

func InitDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=UTC",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("database open: %w", err)
	}

	if cfg.AutoMigrate {
		if err := db.AutoMigrate(
			&model.User{},
			&model.Plan{},
			&model.Order{},
			&model.Payment{},
			&model.InviteCode{},
			&model.Coupon{},
			&model.CommissionLog{},
			&model.Notice{},
			&model.Config{},
			&model.ServerGroup{},
			&model.ServerRoute{},
			&model.ServerTrojan{},
			&model.ServerVmess{},
			&model.ServerShadowsocks{},
			&model.ServerHysteria{},
			&model.ServerNodeSSH{},
			&model.Stat{},
			&model.StatServer{},
			&model.StatUser{},
			&model.Ticket{},
			&model.TicketMessage{},
			&model.Knowledge{},
			&model.Log{},
		); err != nil {
			return nil, fmt.Errorf("database automigrate: %w", err)
		}
	}
	return db, nil
}
