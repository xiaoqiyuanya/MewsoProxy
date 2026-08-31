package admin

import (
	"context"
	"time"

	"mewsoproxy/server/dto"
	"mewsoproxy/server/model"
	"mewsoproxy/server/pkg/apperror"
	redispkg "mewsoproxy/server/pkg/redis"
)

var defaultConfig = map[string]string{
	model.ConfigKeySiteName:       "MewsoProxy",
	model.ConfigKeyRegisterEnabled: "true",
	model.ConfigKeySubscribeURL:   "http://localhost:8081",
	model.ConfigKeySecurePath:     "/admin",
}

func (s *Service) GetConfig(ctx context.Context) ([]dto.AdminConfigDTO, error) {
	var list []model.Config
	if err := s.db.WithContext(ctx).Find(&list).Error; err != nil {
		return nil, apperror.Wrap(apperror.CodeDBError, "配置查询失败", err)
	}
	merged := map[string]string{}
	for k, v := range defaultConfig {
		merged[k] = v
	}
	for _, c := range list {
		merged[c.Key] = c.Value
	}
	out := make([]dto.AdminConfigDTO, 0, len(merged))
	for k, v := range merged {
		out = append(out, dto.AdminConfigDTO{Key: k, Value: v})
	}
	return out, nil
}

func (s *Service) SaveConfig(ctx context.Context, req dto.AdminConfigReq) error {
	now := nowUnix()
	var c model.Config
	err := s.db.WithContext(ctx).Where("key = ?", req.Key).First(&c).Error
	if err == nil {
		return s.db.WithContext(ctx).Model(&model.Config{}).Where("key = ?", req.Key).
			Update("value", req.Value).Update("updated_at", now).Error
	}
	n := model.Config{Key: req.Key, Value: req.Value, CreatedAt: now, UpdatedAt: now}
	return s.db.WithContext(ctx).Create(&n).Error
}

func (s *Service) SystemStatus(ctx context.Context) (dto.AdminSystemStatusDTO, error) {
	res := dto.AdminSystemStatusDTO{ServerTime: time.Now().UTC().Unix(), DBStatus: "ok", RedisStatus: "ok"}
	var userCount, orderCount, todayPaid int64
	_ = s.db.WithContext(ctx).Model(&model.User{}).Count(&userCount)
	_ = s.db.WithContext(ctx).Model(&model.Order{}).Count(&orderCount)
	_ = s.db.WithContext(ctx).Model(&model.Order{}).
		Where("status = ? AND paid_at >= ?", model.OrderStatusCompleted, todayStart()).
		Select("COALESCE(SUM(total_amount),0)").Scan(&todayPaid)
	res.UserCount = userCount
	res.OrderCount = orderCount
	res.TodayPaidTotal = todayPaid
	if s.rds != nil {
		online, _ := s.rds.CountKeys(ctx, redispkg.RedisKeyUserToken+"*")
		res.OnlineUserCount = online
	}
	day := todayRecordAt()
	var onlineNodeCount, activeUserCount, todayTraffic int64
	_ = s.db.WithContext(ctx).Model(&model.ServerNodeTelemetry{}).Where("online = ?", true).Count(&onlineNodeCount)
	_ = s.db.WithContext(ctx).Model(&model.StatUser{}).Where("record_at = ?", day).Distinct("user_id").Count(&activeUserCount)
	_ = s.db.WithContext(ctx).Model(&model.StatServer{}).Where("record_at = ?", day).
		Select("COALESCE(SUM(u+d),0)").Scan(&todayTraffic)
	res.OnlineNodeCount = onlineNodeCount
	res.ActiveUserCount = activeUserCount
	res.TodayTraffic = todayTraffic
	return res, nil
}

func todayStart() int64 {
	now := time.Now().UTC()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return start.Unix()
}

func todayRecordAt() int {
	now := time.Now().UTC()
	return now.Year()*10000 + int(now.Month())*100 + now.Day()
}
