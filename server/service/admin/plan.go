package admin

import (
	"context"
	"time"

	"mewsoproxy/server/dto"
	"mewsoproxy/server/model"
	"mewsoproxy/server/pkg/apperror"
)

func (s *Service) ListPlans(ctx context.Context) ([]model.Plan, error) {
	var list []model.Plan
	if err := s.db.WithContext(ctx).Order("sort asc").Find(&list).Error; err != nil {
		return nil, apperror.Wrap(apperror.CodeDBError, "套餐查询失败", err)
	}
	return list, nil
}

func (s *Service) SavePlan(ctx context.Context, req dto.AdminPlanSaveReq) (int, error) {
	now := time.Now().UTC().Unix()
	p := model.Plan{
		ID:                 req.ID,
		GroupID:            req.GroupID,
		TransferEnable:     req.TransferEnable,
		Name:               req.Name,
		SpeedLimit:         req.SpeedLimit,
		Show:               req.Show,
		Sort:               req.Sort,
		Renew:              req.Renew,
		Content:            req.Content,
		MonthPrice:         req.MonthPrice,
		QuarterPrice:       req.QuarterPrice,
		HalfYearPrice:      req.HalfYearPrice,
		YearPrice:          req.YearPrice,
		TwoYearPrice:       req.TwoYearPrice,
		ThreeYearPrice:     req.ThreeYearPrice,
		OnetimePrice:       req.OnetimePrice,
		ResetPrice:         req.ResetPrice,
		ResetTrafficMethod: req.ResetTrafficMethod,
		CapacityLimit:      req.CapacityLimit,
		UpdatedAt:          now,
	}
	if req.ID > 0 {
		if err := s.db.WithContext(ctx).Model(&model.Plan{}).Where("id = ?", req.ID).Updates(&p).Error; err != nil {
			return 0, apperror.Wrap(apperror.CodeDBError, "保存套餐失败", err)
		}
		return req.ID, nil
	}
	p.CreatedAt = now
	if err := s.db.WithContext(ctx).Create(&p).Error; err != nil {
		return 0, apperror.Wrap(apperror.CodeDBError, "创建套餐失败", err)
	}
	return p.ID, nil
}

func (s *Service) DropPlan(ctx context.Context, id int) error {
	res := s.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Plan{})
	if res.Error != nil {
		return apperror.Wrap(apperror.CodeDBError, "删除套餐失败", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.New(apperror.CodeResourceNotFnd, "套餐不存在")
	}
	return nil
}
