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
	ID               uint    `json:"id" gorm:"primarykey;column:id"`
	Code             string  `json:"code" gorm:"column:code"`
	Name             string  `json:"name" gorm:"column:name"`
	Type             int8    `json:"type" gorm:"column:type"`
	Value            int     `json:"value" gorm:"column:value"`
	Show             bool    `json:"show" gorm:"column:show"`
	LimitUse         *int    `json:"limit_use" gorm:"column:limit_use"`
	LimitUseWithUser *int    `json:"limit_use_with_user" gorm:"column:limit_use_with_user"`
	LimitPlanIDs     *string `json:"limit_plan_ids" gorm:"column:limit_plan_ids"`
	LimitPeriod      *string `json:"limit_period" gorm:"column:limit_period"`
	StartedAt        int64   `json:"started_at" gorm:"column:started_at"`
	EndedAt          int64   `json:"ended_at" gorm:"column:ended_at"`
	CreatedAt        int64   `json:"created_at" gorm:"column:created_at"`
	UpdatedAt        int64   `json:"updated_at" gorm:"column:updated_at"`
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
	ID        uint    `json:"id" gorm:"primarykey;column:id"`
	Title     string  `json:"title" gorm:"column:title"`
	Content   string  `json:"content" gorm:"column:content"`
	Show      bool    `json:"show" gorm:"column:show"`
	ImgURL    *string `json:"img_url" gorm:"column:img_url"`
	Tags      *string `json:"tags" gorm:"column:tags"`
	CreatedAt int64   `json:"created_at" gorm:"column:created_at"`
	UpdatedAt int64   `json:"updated_at" gorm:"column:updated_at"`
}

func (Notice) TableName() string { return "v2_notice" }
