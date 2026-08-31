package model

type Payment struct {
	ID                 int      `json:"id" gorm:"primarykey;column:id"`
	UUID               string   `json:"uuid" gorm:"column:uuid"`
	Payment            string   `json:"payment" gorm:"column:payment"`
	Name               string   `json:"name" gorm:"column:name"`
	Icon               *string  `json:"icon" gorm:"column:icon"`
	Config             string   `json:"config" gorm:"column:config"`
	NotifyDomain       *string  `json:"notify_domain" gorm:"column:notify_domain"`
	HandlingFeeFixed   *int     `json:"handling_fee_fixed" gorm:"column:handling_fee_fixed"`
	HandlingFeePercent *float64 `json:"handling_fee_percent" gorm:"column:handling_fee_percent"`
	Enable             bool     `json:"enable" gorm:"column:enable"`
	Sort               *int     `json:"sort" gorm:"column:sort"`
	CreatedAt          int64    `json:"created_at" gorm:"column:created_at"`
	UpdatedAt          int64    `json:"updated_at" gorm:"column:updated_at"`
}

func (Payment) TableName() string { return "v2_payment" }
