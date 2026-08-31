package database

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"mewsoproxy/server/config"
	"mewsoproxy/server/model"
	"mewsoproxy/server/pkg/token"
)

func SeedAdmin(ctx context.Context, db *gorm.DB, cfg *config.Config) error {
	var cnt int64
	if err := db.WithContext(ctx).Model(&model.User{}).Where("is_admin = ?", true).Count(&cnt).Error; err != nil {
		return fmt.Errorf("count admin: %w", err)
	}
	if cnt > 0 {
		return nil
	}
	email := cfg.App.AdminEmail
	password := cfg.App.AdminPassword
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}
	uuid, _ := token.RandomString(16)
	tk, _ := token.RandomString(16)
	groupID := cfg.App.DefaultGroupID
	now := time.Now().UTC().Unix()
	admin := &model.User{
		Email:        email,
		Password:     string(hash),
		UUID:         uuid,
		Token:        tk,
		GroupID:      &groupID,
		IsAdmin:      true,
		Banned:       false,
		CreatedAt:    now,
		UpdatedAt:    now,
		ExpiredAt:    0,
		TransferEnable: 0,
	}
	if err := db.WithContext(ctx).Create(admin).Error; err != nil {
		return fmt.Errorf("seed admin: %w", err)
	}
	slog.Warn("默认管理员已创建，请尽快在后台修改密码", "email", email)
	return nil
}

func ApplySystemConfig(ctx context.Context, db *gorm.DB, cfg *config.Config) error {
	var rows []model.Config
	if err := db.WithContext(ctx).
		Where("`key` IN ?", []string{model.ConfigKeySubscribeURL, model.ConfigKeyRegisterEnabled}).
		Find(&rows).Error; err != nil {
		return fmt.Errorf("load system config: %w", err)
	}
	for _, row := range rows {
		switch row.Key {
		case model.ConfigKeySubscribeURL:
			if row.Value != "" {
				cfg.App.SubscribeURL = row.Value
			}
		case model.ConfigKeyRegisterEnabled:
			if b, err := strconv.ParseBool(row.Value); err == nil {
				cfg.App.RegisterEnabled = b
			}
		}
	}
	return nil
}
