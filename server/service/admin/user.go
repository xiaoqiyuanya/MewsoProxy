package admin

import (
	"context"
	"time"

	"mewsoproxy/server/dto"
	"mewsoproxy/server/model"
	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/token"
)

func (s *Service) ListUsers(ctx context.Context, keyword string, page, size int) ([]model.User, int64, error) {
	q := s.db.WithContext(ctx).Model(&model.User{})
	if keyword != "" {
		q = q.Where("email LIKE ? OR id = ?", "%"+keyword+"%", keyword)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, apperror.Wrap(apperror.CodeDBError, "用户查询失败", err)
	}
	var list []model.User
	if err := q.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, apperror.Wrap(apperror.CodeDBError, "用户查询失败", err)
	}
	return list, total, nil
}

func (s *Service) GetUser(ctx context.Context, id uint) (*model.User, error) {
	var u model.User
	if err := s.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, apperror.New(apperror.CodeResourceNotFnd, "用户不存在")
	}
	return &u, nil
}

func (s *Service) UpdateUser(ctx context.Context, req dto.AdminUserUpdateReq) error {
	updates := map[string]interface{}{
		"updated_at": time.Now().UTC().Unix(),
	}
	if req.Balance != nil {
		updates["balance"] = *req.Balance
	}
	if req.GroupID != nil {
		updates["group_id"] = *req.GroupID
	}
	if req.PlanID != nil {
		updates["plan_id"] = *req.PlanID
	}
	if req.ExpiredAt != nil {
		updates["expired_at"] = *req.ExpiredAt
	}
	if req.Banned != nil {
		updates["banned"] = *req.Banned
	}
	res := s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", req.ID).Updates(updates)
	if res.Error != nil {
		return apperror.Wrap(apperror.CodeDBError, "更新用户失败", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.New(apperror.CodeResourceNotFnd, "用户不存在")
	}
	return nil
}

func (s *Service) BanUser(ctx context.Context, id uint, ban bool) error {
	res := s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).
		Update("banned", ban).Update("updated_at", time.Now().UTC().Unix())
	if res.Error != nil {
		return apperror.Wrap(apperror.CodeDBError, "操作失败", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.New(apperror.CodeResourceNotFnd, "用户不存在")
	}
	return nil
}

func (s *Service) ResetSecret(ctx context.Context, id uint) error {
	uuid, _ := token.RandomString(16)
	secret, _ := token.RandomString(16)
	res := s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).
		Update("uuid", uuid).Update("token", secret).Update("updated_at", time.Now().UTC().Unix())
	if res.Error != nil {
		return apperror.Wrap(apperror.CodeDBError, "重置失败", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.New(apperror.CodeResourceNotFnd, "用户不存在")
	}
	return nil
}
