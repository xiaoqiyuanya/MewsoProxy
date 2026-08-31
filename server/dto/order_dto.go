package dto

import "mewsoproxy/server/model"

type CreateOrderReq struct {
	PlanID     int    `json:"plan_id" binding:"required,min=1"`
	Period     string `json:"period" binding:"required"`
	CouponCode string `json:"coupon_code" binding:"omitempty"`
	PaymentID  int    `json:"payment_id" binding:"omitempty,min=1"`
	UseBalance bool   `json:"use_balance"`
}

type OrderDTO struct {
	ID           uint   `json:"id"`
	PlanID       int    `json:"plan_id"`
	Period       string `json:"period"`
	TradeNo      string `json:"trade_no"`
	Type         int    `json:"type"`
	TotalAmount  int    `json:"total_amount"`
	BalanceAmount int   `json:"balance_amount"`
	DiscountAmount int  `json:"discount_amount"`
	Status       int8   `json:"status"`
	CreatedAt    Time   `json:"created_at"`
	PaidAt       Time   `json:"paid_at"`
}

func ToOrderDTO(o *model.Order) OrderDTO {
	return OrderDTO{
		ID:            o.ID,
		PlanID:        o.PlanID,
		Period:        o.Period,
		TradeNo:       o.TradeNo,
		Type:          o.Type,
		TotalAmount:   o.TotalAmount,
		BalanceAmount: orZero(o.BalanceAmount),
		DiscountAmount: orZero(o.DiscountAmount),
		Status:        o.Status,
		CreatedAt:     FromUnix(o.CreatedAt),
		PaidAt:        fromUnixPtr(o.PaidAt),
	}
}

func orZero(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func fromUnixPtr(p *int64) Time {
	if p == nil {
		return Time{}
	}
	return FromUnix(*p)
}
