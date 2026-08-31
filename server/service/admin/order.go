package admin

import (
	"context"

	"mewsoproxy/server/model"
	"mewsoproxy/server/pkg/apperror"
)

func (s *Service) ListOrders(ctx context.Context, status int8, keyword string, page, size int) ([]model.Order, int64, error) {
	q := s.db.WithContext(ctx).Model(&model.Order{})
	if status >= 0 {
		q = q.Where("status = ?", status)
	}
	if keyword != "" {
		q = q.Where("trade_no LIKE ? OR user_id = ?", "%"+keyword+"%", keyword)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, apperror.Wrap(apperror.CodeDBError, "订单查询失败", err)
	}
	var list []model.Order
	if err := q.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, apperror.Wrap(apperror.CodeDBError, "订单查询失败", err)
	}
	return list, total, nil
}

func (s *Service) CancelOrder(ctx context.Context, orderID uint) error {
	var o model.Order
	if err := s.db.WithContext(ctx).First(&o, orderID).Error; err != nil {
		return apperror.New(apperror.CodeResourceNotFnd, "订单不存在")
	}
	if o.Status != model.OrderStatusPending {
		return apperror.New(apperror.CodeOrderPaid, "订单当前不可取消")
	}
	return s.orderSvc.Cancel(ctx, o.UserID, o.ID)
}

func (s *Service) MarkOrderPaid(ctx context.Context, orderID uint, callbackNo string) error {
	return s.orderSvc.MarkPaid(ctx, orderID, callbackNo)
}

func (s *Service) GetOrder(ctx context.Context, orderID uint) (*model.Order, error) {
	var o model.Order
	if err := s.db.WithContext(ctx).First(&o, orderID).Error; err != nil {
		return nil, apperror.New(apperror.CodeResourceNotFnd, "订单不存在")
	}
	return &o, nil
}
