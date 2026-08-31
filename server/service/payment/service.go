package payment

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"mewsoproxy/server/config"
	"mewsoproxy/server/model"
	"mewsoproxy/server/pkg/apperror"
	ordersvc "mewsoproxy/server/service/order"
)

type Service struct {
	db       *gorm.DB
	cfg      *config.Config
	orderSvc *ordersvc.Service
}

func NewService(db *gorm.DB, cfg *config.Config, orderSvc *ordersvc.Service) *Service {
	return &Service{db: db, cfg: cfg, orderSvc: orderSvc}
}

func (s *Service) Create(ctx context.Context, userID uint, orderID uint, paymentID int, baseURL string) (*CreateResult, error) {
	var o model.Order
	if err := s.db.WithContext(ctx).First(&o, orderID).Error; err != nil {
		return nil, apperror.New(apperror.CodeResourceNotFnd, "订单不存在")
	}
	if o.UserID != userID {
		return nil, apperror.New(apperror.CodeNoPermission, "无权操作该订单")
	}
	if o.Status != model.OrderStatusPending {
		return nil, apperror.New(apperror.CodeOrderPaid, "订单当前不可支付")
	}
	payable := o.TotalAmount - ptrValue(o.BalanceAmount)
	if payable <= 0 {
		if err := s.orderSvc.MarkPaidByTradeNo(ctx, o.TradeNo, "BALANCE"); err != nil {
			return nil, err
		}
		return &CreateResult{PayType: "completed", Completed: true}, nil
	}
	p, err := s.resolvePayment(ctx, orderID, paymentID)
	if err != nil {
		return nil, err
	}
	if o.PaymentID == nil || *o.PaymentID != p.ID {
		now := time.Now().UTC().Unix()
		if err := s.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", o.ID).
			Update("payment_id", p.ID).Update("updated_at", now).Error; err != nil {
			return nil, apperror.Wrap(apperror.CodeDBError, "保存支付渠道失败", err)
		}
	}
	domain := baseURL
	if p.NotifyDomain != nil && *p.NotifyDomain != "" {
		domain = *p.NotifyDomain
	}
	gw, err := Resolve(p)
	if err != nil {
		return nil, err
	}
	in := CreateInput{
		Order:     &o,
		Payment:   p,
		NotifyURL: fmt.Sprintf("%s/api/v1/payment/notify", domain),
		ReturnURL: fmt.Sprintf("%s/#/order", domain),
	}
	return gw.Create(ctx, in)
}

func (s *Service) Channels(ctx context.Context) ([]ChannelDTO, error) {
	var list []model.Payment
	if err := s.db.WithContext(ctx).Where("enable = ?", true).Order("sort asc").Find(&list).Error; err != nil {
		return nil, apperror.Wrap(apperror.CodeDBError, "支付渠道查询失败", err)
	}
	out := make([]ChannelDTO, 0, len(list))
	for i := range list {
		out = append(out, ChannelDTO{
			ID:     list[i].ID,
			Payment: list[i].Payment,
			Name:   list[i].Name,
			Icon:   list[i].Icon,
		})
	}
	return out, nil
}

func (s *Service) Notify(ctx context.Context, params map[string]string) error {
	tradeNo := params["out_trade_no"]
	if tradeNo == "" {
		tradeNo = params["trade_no"]
	}
	if tradeNo == "" {
		return apperror.New(apperror.CodeParamFormat, "订单号缺失")
	}
	var o model.Order
	if err := s.db.WithContext(ctx).Where("trade_no = ?", tradeNo).First(&o).Error; err != nil {
		return apperror.New(apperror.CodeResourceNotFnd, "订单不存在")
	}
	callbackNo := params["trade_no"]
	if callbackNo == "" {
		callbackNo = params["callback_no"]
	}
	if o.PaymentID == nil {
		return s.orderSvc.MarkPaidByTradeNo(ctx, tradeNo, callbackNo)
	}
	var p model.Payment
	if err := s.db.WithContext(ctx).First(&p, *o.PaymentID).Error; err != nil {
		return apperror.New(apperror.CodeResourceNotFnd, "支付渠道不存在")
	}
	gw, err := Resolve(&p)
	if err != nil {
		return err
	}
	if err := gw.VerifyNotify(ctx, params); err != nil {
		return err
	}
	return s.orderSvc.MarkPaidByTradeNo(ctx, tradeNo, callbackNo)
}

func (s *Service) resolvePayment(ctx context.Context, orderID uint, paymentID int) (*model.Payment, error) {
	id := paymentID
	if id == 0 {
		var o model.Order
		if err := s.db.WithContext(ctx).Select("payment_id").First(&o, orderID).Error; err == nil && o.PaymentID != nil {
			id = *o.PaymentID
		}
	}
	if id > 0 {
		var p model.Payment
		if err := s.db.WithContext(ctx).First(&p, id).Error; err != nil {
			return nil, apperror.New(apperror.CodeResourceNotFnd, "支付渠道不存在")
		}
		if !p.Enable {
			return nil, apperror.New(apperror.CodeParamFormat, "支付渠道未启用")
		}
		return &p, nil
	}
	var p model.Payment
	if err := s.db.WithContext(ctx).Where("enable = ?", true).Order("sort asc").First(&p).Error; err != nil {
		return nil, apperror.New(apperror.CodeParamFormat, "暂无可用支付渠道")
	}
	return &p, nil
}

func ptrValue(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}
