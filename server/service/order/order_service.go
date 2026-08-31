package order

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"mewsoproxy/server/config"
	"mewsoproxy/server/dto"
	"mewsoproxy/server/model"
	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/token"
)

type Service struct {
	db  *gorm.DB
	cfg *config.Config
}

func New(db *gorm.DB, cfg *config.Config) *Service {
	return &Service{db: db, cfg: cfg}
}

var periodDuration = map[string]int64{
	"month_price":      30 * 86400,
	"quarter_price":    90 * 86400,
	"half_year_price":  180 * 86400,
	"year_price":       365 * 86400,
	"two_year_price":   730 * 86400,
	"three_year_price": 1095 * 86400,
	"onetime_price":    30 * 86400,
}

type CreateOrderResult struct {
	Order *model.Order
	Total int
}

func (s *Service) Create(ctx context.Context, userID uint, req dto.CreateOrderReq) (*model.Order, error) {
	plan, err := s.getPlan(ctx, req.PlanID)
	if err != nil {
		return nil, err
	}
	price, _, ok := periodPrice(plan, req.Period)
	if !ok {
		return nil, apperror.New(apperror.CodeParamFormat, "订单周期不合法")
	}
	now := time.Now().UTC()
	user, err := s.getUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	otype := model.OrderTypeNew
	if user.PlanID != nil {
		if *user.PlanID == req.PlanID && user.ExpiredAt > 0 {
			otype = model.OrderTypeRenew
		} else if *user.PlanID != req.PlanID {
			otype = model.OrderTypeUpgrade
		}
	}
	discount := 0
	if otype == model.OrderTypeUpgrade {
		surplus, ok := s.surplusValue(ctx, user, plan)
		if !ok {
			return nil, apperror.New(apperror.CodeOrderInvalid, "无法计算剩余价值")
		}
		if surplus > price {
			surplus = price
		}
		discount = surplus
	}
	couponID, couponDiscount, err := s.applyCoupon(ctx, userID, req.CouponCode, req.PlanID, req.Period)
	if err != nil {
		return nil, err
	}
	discount += couponDiscount

	total := price - discount
	if total < 0 {
		total = 0
	}
	finalTotal := total

	var balanceAmount int
	if req.UseBalance {
		if user.Balance <= 0 {
			return nil, apperror.New(apperror.CodeBalanceNotEnoug, "余额不足")
		}
		use := min(user.Balance, finalTotal)
		balanceAmount = use
		finalTotal = finalTotal - use
	}

	tradeNo := newTradeNo()
	o := &model.Order{
		UserID:          userID,
		PlanID:          req.PlanID,
		Type:            otype,
		Period:          req.Period,
		TradeNo:         tradeNo,
		TotalAmount:     total,
		BalanceAmount:   &balanceAmount,
		DiscountAmount:  &discount,
		CouponID:        couponID,
		Status:          model.OrderStatusPending,
		CommissionStatus: 0,
		CommissionBalance: 0,
		CreatedAt:       now.Unix(),
		UpdatedAt:       now.Unix(),
	}
	if user.InviteUserID != nil {
		o.InviteUserID = user.InviteUserID
	}
	if req.PaymentID > 0 {
		o.PaymentID = &req.PaymentID
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if balanceAmount > 0 {
			if err := tx.Model(&model.User{}).
				Where("id = ? AND balance >= ?", userID, balanceAmount).
				UpdateColumn("balance", gorm.Expr("balance - ?", balanceAmount)).Error; err != nil {
				return err
			}
		}
		return tx.Create(o).Error
	})
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeDBError, "创建订单失败", err)
	}
	return o, nil
}

func (s *Service) List(ctx context.Context, userID uint, page, size int) ([]model.Order, int64, error) {
	var total int64
	if err := s.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, apperror.Wrap(apperror.CodeDBError, "订单查询失败", err)
	}
	var list []model.Order
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, apperror.Wrap(apperror.CodeDBError, "订单查询失败", err)
	}
	return list, total, nil
}

func (s *Service) Detail(ctx context.Context, userID, orderID uint) (*model.Order, error) {
	var o model.Order
	if err := s.db.WithContext(ctx).First(&o, orderID).Error; err != nil {
		return nil, apperror.New(apperror.CodeResourceNotFnd, "订单不存在")
	}
	if o.UserID != userID {
		return nil, apperror.New(apperror.CodeNoPermission, "无权访问该订单")
	}
	return &o, nil
}

func (s *Service) Cancel(ctx context.Context, userID, orderID uint) error {
	var o model.Order
	if err := s.db.WithContext(ctx).First(&o, orderID).Error; err != nil {
		return apperror.New(apperror.CodeResourceNotFnd, "订单不存在")
	}
	if o.UserID != userID {
		return apperror.New(apperror.CodeNoPermission, "无权操作该订单")
	}
	if o.Status != model.OrderStatusPending {
		return apperror.New(apperror.CodeOrderPaid, "订单当前不可取消")
	}
	o.Status = model.OrderStatusCanceled
	o.UpdatedAt = time.Now().UTC().Unix()
	if o.BalanceAmount != nil && *o.BalanceAmount > 0 {
		if err := s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).
			UpdateColumn("balance", gorm.Expr("balance + ?", *o.BalanceAmount)).Error; err != nil {
			return apperror.Wrap(apperror.CodeDBError, "退款失败", err)
		}
	}
	if err := s.db.WithContext(ctx).Model(&o).Update("status", o.Status).Update("updated_at", o.UpdatedAt).Error; err != nil {
		return apperror.Wrap(apperror.CodeDBError, "取消失败", err)
	}
	return nil
}

func (s *Service) MarkPaid(ctx context.Context, orderID uint, callbackNo string) error {
	var o model.Order
	if err := s.db.WithContext(ctx).First(&o, orderID).Error; err != nil {
		return apperror.New(apperror.CodeResourceNotFnd, "订单不存在")
	}
	if o.Status != model.OrderStatusPending {
		return nil
	}
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, o.UserID).Error; err != nil {
		return apperror.Wrap(apperror.CodeDBError, "用户不存在", err)
	}
	var plan model.Plan
	if err := s.db.WithContext(ctx).First(&plan, o.PlanID).Error; err != nil {
		return apperror.New(apperror.CodeResourceNotFnd, "套餐不存在")
	}
	_, duration, ok := periodPrice(&plan, o.Period)
	if !ok {
		return apperror.New(apperror.CodeOrderInvalid, "订单周期不合法")
	}
	now := time.Now().UTC().Unix()
	expire := now + duration
	if o.Type != model.OrderTypeNew && user.ExpiredAt > now {
		expire = user.ExpiredAt + duration
	}
	s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.User{}).Where("id = ?", o.UserID).
			Update("expired_at", expire).
			Update("plan_id", o.PlanID).
			Update("transfer_enable", plan.TransferEnable).
			Update("u", 0).
			Update("d", 0).Error; err != nil {
			return err
		}
		cb := callbackNo
		return tx.Model(&o).Where("status = ?", model.OrderStatusPending).
			Update("status", model.OrderStatusCompleted).
			Update("callback_no", cb).
			Update("paid_at", now).
			Update("updated_at", now).Error
	})
	s.grantCommission(ctx, &o)
	return nil
}

func (s *Service) MarkPaidByTradeNo(ctx context.Context, tradeNo, callbackNo string) error {
	var o model.Order
	if err := s.db.WithContext(ctx).Where("trade_no = ?", tradeNo).First(&o).Error; err != nil {
		return apperror.New(apperror.CodeResourceNotFnd, "订单不存在")
	}
	return s.MarkPaid(ctx, o.ID, callbackNo)
}

func (s *Service) grantCommission(ctx context.Context, o *model.Order) {
	if o.InviteUserID == nil {
		return
	}
	var inviter model.User
	if err := s.db.WithContext(ctx).First(&inviter, *o.InviteUserID).Error; err != nil {
		return
	}
	if inviter.CommissionRate == nil || *inviter.CommissionRate <= 0 {
		return
	}
	get := int64(o.TotalAmount) * int64(*inviter.CommissionRate) / 100
	if get <= 0 {
		return
	}
	log := &model.CommissionLog{
		InviteUserID: inviter.ID,
		UserID:       o.UserID,
		TradeNo:      o.TradeNo,
		OrderAmount:  o.TotalAmount,
		GetAmount:    int(get),
		CreatedAt:    time.Now().UTC().Unix(),
		UpdatedAt:    time.Now().UTC().Unix(),
	}
	_ = s.db.WithContext(ctx).Create(log)
	_ = s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", inviter.ID).
		UpdateColumn("commission_balance", gorm.Expr("commission_balance + ?", int(get))).Error
}

func (s *Service) getPlan(ctx context.Context, id int) (*model.Plan, error) {
	var p model.Plan
	if err := s.db.WithContext(ctx).First(&p, id).Error; err != nil {
		return nil, apperror.New(apperror.CodeResourceNotFnd, "套餐不存在")
	}
	return &p, nil
}

func (s *Service) getUser(ctx context.Context, id uint) (*model.User, error) {
	var u model.User
	if err := s.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, apperror.New(apperror.CodeUserNotFound, "用户不存在")
	}
	return &u, nil
}

func (s *Service) surplusValue(ctx context.Context, user *model.User, newPlan *model.Plan) (int, bool) {
	if user.ExpiredAt <= 0 || user.PlanID == nil {
		return 0, true
	}
	var curPlan model.Plan
	if err := s.db.WithContext(ctx).First(&curPlan, *user.PlanID).Error; err != nil {
		return 0, false
	}
	remaining := user.ExpiredAt - time.Now().UTC().Unix()
	if remaining <= 0 {
		return 0, true
	}
	price := maxPeriodPrice(&curPlan)
	if price == 0 {
		return 0, true
	}
	perDay := price / 30
	return int(int64(perDay) * remaining / 86400), true
}

func (s *Service) applyCoupon(ctx context.Context, userID uint, code string, planID int, period string) (*int, int, error) {
	if code == "" {
		return nil, 0, nil
	}
	var c model.Coupon
	if err := s.db.WithContext(ctx).Where("code = ?", code).First(&c).Error; err != nil {
		return nil, 0, apperror.New(apperror.CodeCouponInvalid, "优惠券无效")
	}
	now := time.Now().UTC().Unix()
	if now < c.StartedAt || (c.EndedAt > 0 && now > c.EndedAt) {
		return nil, 0, apperror.New(apperror.CodeCouponInvalid, "优惠券不在有效期内")
	}
	if c.LimitPlanIDs != nil && *c.LimitPlanIDs != "" && !strings.Contains(*c.LimitPlanIDs + ",", fmt.Sprintf("%d,", planID)) {
		return nil, 0, apperror.New(apperror.CodeCouponInvalid, "优惠券不适用于该套餐")
	}
	var used int64
	s.db.WithContext(ctx).Model(&model.Order{}).Where("coupon_id = ?", c.ID).Count(&used)
	if c.LimitUse != nil && int(used) >= *c.LimitUse {
		return nil, 0, apperror.New(apperror.CodeCouponInvalid, "优惠券已达使用上限")
	}
	cid := int(c.ID)
	return &cid, c.Value, nil
}

func periodPrice(p *model.Plan, period string) (int, int64, bool) {
	var price *int
	switch period {
	case "month_price":
		price = p.MonthPrice
	case "quarter_price":
		price = p.QuarterPrice
	case "half_year_price":
		price = p.HalfYearPrice
	case "year_price":
		price = p.YearPrice
	case "two_year_price":
		price = p.TwoYearPrice
	case "three_year_price":
		price = p.ThreeYearPrice
	case "onetime_price":
		price = p.OnetimePrice
	default:
		return 0, 0, false
	}
	if price == nil {
		return 0, 0, false
	}
	dur, ok := periodDuration[period]
	if !ok {
		dur = 30 * 86400
	}
	return *price, dur, true
}

func priceByPeriod(p *model.Plan, period string) int {
	price, _, _ := periodPrice(p, period)
	return price
}

func maxPeriodPrice(p *model.Plan) int {
	m := 0
	for _, period := range []string{"month_price", "quarter_price", "half_year_price", "year_price", "two_year_price", "three_year_price"} {
		if v := priceByPeriod(p, period); v > m {
			m = v
		}
	}
	return m
}

func newTradeNo() string {
	s, _ := token.RandomString(16)
	return "ORD-" + s
}
