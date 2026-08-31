package model

type InviteCode struct {
	ID        uint   `gorm:"primarykey;column:id"`
	UserID    uint   `gorm:"column:user_id"`
	Code      string `gorm:"column:code"`
	Status    bool   `gorm:"column:status"`
	PV        int    `gorm:"column:pv"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

func (InviteCode) TableName() string { return "v2_invite_code" }

type Coupon struct {
	ID               uint   `gorm:"primarykey;column:id"`
	Code             string `gorm:"column:code"`
	Name             string `gorm:"column:name"`
	Type             int8   `gorm:"column:type"`
	Value            int    `gorm:"column:value"`
	Show             bool   `gorm:"column:show"`
	LimitUse         *int   `gorm:"column:limit_use"`
	LimitUseWithUser *int   `gorm:"column:limit_use_with_user"`
	LimitPlanIDs     *string `gorm:"column:limit_plan_ids"`
	LimitPeriod      *string `gorm:"column:limit_period"`
	StartedAt        int64  `gorm:"column:started_at"`
	EndedAt          int64  `gorm:"column:ended_at"`
	CreatedAt        int64  `gorm:"column:created_at"`
	UpdatedAt        int64  `gorm:"column:updated_at"`
}

func (Coupon) TableName() string { return "v2_coupon" }

type CommissionLog struct {
	ID           uint   `gorm:"primarykey;column:id"`
	InviteUserID uint   `gorm:"column:invite_user_id"`
	UserID       uint   `gorm:"column:user_id"`
	TradeNo      string `gorm:"column:trade_no"`
	OrderAmount  int    `gorm:"column:order_amount"`
	GetAmount    int    `gorm:"column:get_amount"`
	CreatedAt    int64  `gorm:"column:created_at"`
	UpdatedAt    int64  `gorm:"column:updated_at"`
}

func (CommissionLog) TableName() string { return "v2_commission_log" }

type Notice struct {
	ID        uint   `gorm:"primarykey;column:id"`
	Title     string `gorm:"column:title"`
	Content   string `gorm:"column:content"`
	Show      bool   `gorm:"column:show"`
	ImgURL    *string `gorm:"column:img_url"`
	Tags      *string `gorm:"column:tags"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

func (Notice) TableName() string { return "v2_notice" }
