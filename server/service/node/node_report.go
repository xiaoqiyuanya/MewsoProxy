package node

import (
	"context"
	"time"

	"gorm.io/gorm"

	"mewsoproxy/server/config"
	"mewsoproxy/server/dto"
	"mewsoproxy/server/model"
	"mewsoproxy/server/pkg/apperror"
)

type Service struct {
	db  *gorm.DB
	cfg *config.Config
}

func New(db *gorm.DB, cfg *config.Config) *Service {
	return &Service{db: db, cfg: cfg}
}

func (s *Service) Report(ctx context.Context, req dto.NodeReportReq) error {
	if s.cfg.App.ServerToken == "" {
		return apperror.New(apperror.CodeNoPermission, "站点未配置 server_token")
	}
	if req.Token != s.cfg.App.ServerToken {
		return apperror.New(apperror.CodeNoPermission, "节点令牌无效")
	}

	if err := s.upsertTelemetry(ctx, req); err != nil {
		return err
	}

	if err := s.aggregateUserTraffic(ctx, req); err != nil {
		return err
	}

	if err := s.aggregateServerTraffic(ctx, req); err != nil {
		return err
	}
	return nil
}

func (s *Service) upsertTelemetry(ctx context.Context, req dto.NodeReportReq) error {
	now := time.Now().UTC().Unix()
	var t model.ServerNodeTelemetry
	online := true
	if req.Online != nil {
		online = *req.Online
	}
	err := s.db.WithContext(ctx).Where("node_type = ? AND node_id = ?", req.NodeType, req.NodeID).First(&t).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return apperror.Wrap(apperror.CodeDBError, "节点遥测查询失败", err)
		}
		t = model.ServerNodeTelemetry{
			NodeType: req.NodeType,
			NodeID:   req.NodeID,
			Online:   online,
			Uptime:   req.Uptime,
			Load:     req.Load,
			CreatedAt: now,
		}
		if online {
			t.LastOnlineAt = now
		}
		t.UpdatedAt = now
		if err := s.db.WithContext(ctx).Create(&t).Error; err != nil {
			return apperror.Wrap(apperror.CodeDBError, "节点遥测写入失败", err)
		}
		return nil
	}

	t.Online = online
	t.Uptime = req.Uptime
	t.Load = req.Load
	if online {
		t.LastOnlineAt = now
	}
	t.UpdatedAt = now
	if err := s.db.WithContext(ctx).Save(&t).Error; err != nil {
		return apperror.Wrap(apperror.CodeDBError, "节点遥测更新失败", err)
	}
	return nil
}

func (s *Service) aggregateUserTraffic(ctx context.Context, req dto.NodeReportReq) error {
	if len(req.Users) == 0 {
		return nil
	}
	day := recordAt()
	for _, item := range req.Users {
		if item.U == 0 && item.D == 0 {
			continue
		}
		var u model.User
		if err := s.db.WithContext(ctx).Where("uuid = ?", item.UUID).First(&u).Error; err != nil {
			continue
		}
		addUserTraffic(ctx, s.db, &u, item.U, item.D)
		upsertStatUser(ctx, s.db, int(u.ID), item.U, item.D, day)
	}
	return nil
}

func (s *Service) aggregateServerTraffic(ctx context.Context, req dto.NodeReportReq) error {
	if req.U == 0 && req.D == 0 {
		return nil
	}
	day := recordAt()
	upsertStatServer(ctx, s.db, req.NodeType, req.NodeID, req.U, req.D, day)
	return nil
}

func addUserTraffic(ctx context.Context, db *gorm.DB, u *model.User, uBytes, dBytes int64) {
	now := time.Now().UTC().Unix()
	db.WithContext(ctx).Model(&model.User{}).Where("id = ?", u.ID).Updates(map[string]interface{}{
		"u":          gorm.Expr("u + ?", uBytes),
		"d":          gorm.Expr("d + ?", dBytes),
		"updated_at": now,
	})
}

func upsertStatUser(ctx context.Context, db *gorm.DB, userID int, uBytes, dBytes int64, day int) {
	var st model.StatUser
	err := db.WithContext(ctx).Where("user_id = ? AND record_at = ?", userID, day).First(&st).Error
	now := time.Now().UTC().Unix()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return
		}
		st = model.StatUser{UserID: userID, U: uBytes, D: dBytes, RecordAt: day, CreatedAt: now, UpdatedAt: now}
		db.WithContext(ctx).Create(&st)
		return
	}
	st.U += uBytes
	st.D += dBytes
	st.UpdatedAt = now
	db.WithContext(ctx).Model(&st).Where("id = ?", st.ID).Updates(map[string]interface{}{
		"u":          st.U,
		"d":          st.D,
		"updated_at": now,
	})
}

func upsertStatServer(ctx context.Context, db *gorm.DB, nodeType string, nodeID int, uBytes, dBytes int64, day int) {
	var st model.StatServer
	err := db.WithContext(ctx).Where("server_id = ? AND server_type = ? AND record_at = ?", nodeID, nodeType, day).First(&st).Error
	now := time.Now().UTC().Unix()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return
		}
		st = model.StatServer{ServerID: nodeID, ServerType: nodeType, U: uBytes, D: dBytes, RecordAt: day, CreatedAt: now, UpdatedAt: now}
		db.WithContext(ctx).Create(&st)
		return
	}
	st.U += uBytes
	st.D += dBytes
	st.UpdatedAt = now
	db.WithContext(ctx).Model(&st).Where("id = ?", st.ID).Updates(map[string]interface{}{
		"u":          st.U,
		"d":          st.D,
		"updated_at": now,
	})
}

func recordAt() int {
	t := time.Now().UTC()
	return t.Year()*10000 + int(t.Month())*100 + t.Day()
}
