package admin

import (
	"context"
	"time"

	"mewsoproxy/server/dto"
	"mewsoproxy/server/model"
	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/token"
)

func nowUnix() int64 { return time.Now().UTC().Unix() }

func (s *Service) ListCoupons(ctx context.Context) ([]model.Coupon, error) {
	var list []model.Coupon
	if err := s.db.WithContext(ctx).Order("id desc").Find(&list).Error; err != nil {
		return nil, apperror.Wrap(apperror.CodeDBError, "优惠券查询失败", err)
	}
	return list, nil
}

func (s *Service) SaveCoupon(ctx context.Context, req dto.AdminCouponSaveReq) (uint, error) {
	c := model.Coupon{
		Code: req.Code, Name: req.Name, Type: req.Type, Value: req.Value,
		Show: req.Show, LimitUse: req.LimitUse, LimitUseWithUser: req.LimitUseWithUser,
		LimitPlanIDs: req.LimitPlanIDs, LimitPeriod: req.LimitPeriod,
		StartedAt: req.StartedAt, EndedAt: req.EndedAt, UpdatedAt: nowUnix(),
	}
	if req.ID > 0 {
		c.ID = req.ID
		if err := s.db.WithContext(ctx).Model(&model.Coupon{}).Where("id = ?", req.ID).Updates(&c).Error; err != nil {
			return 0, apperror.Wrap(apperror.CodeDBError, "保存优惠券失败", err)
		}
		return req.ID, nil
	}
	c.CreatedAt = nowUnix()
	if err := s.db.WithContext(ctx).Create(&c).Error; err != nil {
		return 0, apperror.Wrap(apperror.CodeDBError, "创建优惠券失败", err)
	}
	return c.ID, nil
}

func (s *Service) DropCoupon(ctx context.Context, id uint) error {
	res := s.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Coupon{})
	if res.Error != nil {
		return apperror.Wrap(apperror.CodeDBError, "删除优惠券失败", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.New(apperror.CodeResourceNotFnd, "优惠券不存在")
	}
	return nil
}

func (s *Service) ToggleCouponShow(ctx context.Context, id uint, show bool) error {
	res := s.db.WithContext(ctx).Model(&model.Coupon{}).Where("id = ?", id).
		Update("show", show).Update("updated_at", nowUnix())
	if res.Error != nil {
		return apperror.Wrap(apperror.CodeDBError, "操作失败", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.New(apperror.CodeResourceNotFnd, "优惠券不存在")
	}
	return nil
}

func (s *Service) ListNotices(ctx context.Context) ([]model.Notice, error) {
	var list []model.Notice
	if err := s.db.WithContext(ctx).Order("id desc").Find(&list).Error; err != nil {
		return nil, apperror.Wrap(apperror.CodeDBError, "公告查询失败", err)
	}
	return list, nil
}

func (s *Service) SaveNotice(ctx context.Context, req dto.AdminNoticeSaveReq) (uint, error) {
	n := model.Notice{
		Title: req.Title, Content: req.Content, Show: req.Show,
		ImgURL: req.ImgURL, Tags: req.Tags, UpdatedAt: nowUnix(),
	}
	if req.ID > 0 {
		n.ID = req.ID
		if err := s.db.WithContext(ctx).Model(&model.Notice{}).Where("id = ?", req.ID).Updates(&n).Error; err != nil {
			return 0, apperror.Wrap(apperror.CodeDBError, "保存公告失败", err)
		}
		return req.ID, nil
	}
	n.CreatedAt = nowUnix()
	if err := s.db.WithContext(ctx).Create(&n).Error; err != nil {
		return 0, apperror.Wrap(apperror.CodeDBError, "创建公告失败", err)
	}
	return n.ID, nil
}

func (s *Service) DropNotice(ctx context.Context, id uint) error {
	res := s.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Notice{})
	if res.Error != nil {
		return apperror.Wrap(apperror.CodeDBError, "删除公告失败", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.New(apperror.CodeResourceNotFnd, "公告不存在")
	}
	return nil
}

func (s *Service) ToggleNoticeShow(ctx context.Context, id uint, show bool) error {
	res := s.db.WithContext(ctx).Model(&model.Notice{}).Where("id = ?", id).
		Update("show", show).Update("updated_at", nowUnix())
	if res.Error != nil {
		return apperror.Wrap(apperror.CodeDBError, "操作失败", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.New(apperror.CodeResourceNotFnd, "公告不存在")
	}
	return nil
}

func (s *Service) ListPayments(ctx context.Context) ([]model.Payment, error) {
	var list []model.Payment
	if err := s.db.WithContext(ctx).Order("sort asc").Find(&list).Error; err != nil {
		return nil, apperror.Wrap(apperror.CodeDBError, "支付渠道查询失败", err)
	}
	return list, nil
}

func (s *Service) SavePayment(ctx context.Context, req dto.AdminPaymentSaveReq) (int, error) {
	p := model.Payment{
		UUID: req.UUID, Payment: req.Payment, Name: req.Name, Icon: req.Icon,
		Config: req.Config, NotifyDomain: req.NotifyDomain,
		HandlingFeeFixed: req.HandlingFeeFixed, HandlingFeePercent: req.HandlingFeePercent,
		Enable: req.Enable, Sort: req.Sort, UpdatedAt: nowUnix(),
	}
	if req.ID > 0 {
		p.ID = req.ID
		if err := s.db.WithContext(ctx).Model(&model.Payment{}).Where("id = ?", req.ID).Updates(&p).Error; err != nil {
			return 0, apperror.Wrap(apperror.CodeDBError, "保存支付渠道失败", err)
		}
		return req.ID, nil
	}
	p.CreatedAt = nowUnix()
	if err := s.db.WithContext(ctx).Create(&p).Error; err != nil {
		return 0, apperror.Wrap(apperror.CodeDBError, "创建支付渠道失败", err)
	}
	return p.ID, nil
}

func (s *Service) DropPayment(ctx context.Context, id int) error {
	res := s.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Payment{})
	if res.Error != nil {
		return apperror.Wrap(apperror.CodeDBError, "删除支付渠道失败", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.New(apperror.CodeResourceNotFnd, "支付渠道不存在")
	}
	return nil
}

func (s *Service) TogglePaymentShow(ctx context.Context, id int, enable bool) error {
	res := s.db.WithContext(ctx).Model(&model.Payment{}).Where("id = ?", id).
		Update("enable", enable).Update("updated_at", nowUnix())
	if res.Error != nil {
		return apperror.Wrap(apperror.CodeDBError, "操作失败", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.New(apperror.CodeResourceNotFnd, "支付渠道不存在")
	}
	return nil
}

func (s *Service) GenerateCouponCode() string {
	v, _ := token.RandomString(6)
	return "CPN-" + v
}
