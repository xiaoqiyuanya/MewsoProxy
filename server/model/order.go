package model

type Order struct {
	ID                    uint   `gorm:"primarykey;column:id"`
	InviteUserID          *uint  `gorm:"column:invite_user_id"`
	UserID                uint   `gorm:"column:user_id;index"`
	PlanID                int    `gorm:"column:plan_id"`
	CouponID              *int   `gorm:"column:coupon_id"`
	PaymentID             *int   `gorm:"column:payment_id"`
	Type                  int    `gorm:"column:type"`
	Period                string `gorm:"column:period"`
	TradeNo               string `gorm:"column:trade_no;type:varchar(36);uniqueIndex"`
	CallbackNo            *string `gorm:"column:callback_no"`
	TotalAmount           int    `gorm:"column:total_amount"`
	HandlingAmount        *int   `gorm:"column:handling_amount"`
	DiscountAmount        *int   `gorm:"column:discount_amount"`
	SurplusAmount         *int   `gorm:"column:surplus_amount"`
	RefundAmount          *int   `gorm:"column:refund_amount"`
	BalanceAmount         *int   `gorm:"column:balance_amount"`
	SurplusOrderIDs       *string `gorm:"column:surplus_order_ids"`
	Status                int8   `gorm:"column:status"`
	CommissionStatus      int8   `gorm:"column:commission_status"`
	CommissionBalance     int    `gorm:"column:commission_balance"`
	ActualCommissionBalance *int `gorm:"column:actual_commission_balance"`
	PaidAt                *int64 `gorm:"column:paid_at"`
	CreatedAt             int64  `gorm:"column:created_at"`
	UpdatedAt             int64  `gorm:"column:updated_at"`
}

func (Order) TableName() string { return "v2_order" }

const (
	OrderTypeNew     = 1
	OrderTypeRenew   = 2
	OrderTypeUpgrade = 3
)

const (
	OrderStatusPending      = 0
	OrderStatusProcessing   = 1
	OrderStatusCanceled     = 2
	OrderStatusCompleted    = 3
	OrderStatusDiscounted   = 4
)
