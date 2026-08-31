package model

type Payment struct {
	ID               int     `gorm:"primarykey;column:id"`
	UUID             string  `gorm:"column:uuid"`
	Payment          string  `gorm:"column:payment"`
	Name             string  `gorm:"column:name"`
	Icon             *string `gorm:"column:icon"`
	Config           string  `gorm:"column:config"`
	NotifyDomain     *string `gorm:"column:notify_domain"`
	HandlingFeeFixed *int    `gorm:"column:handling_fee_fixed"`
	HandlingFeePercent *float64 `gorm:"column:handling_fee_percent"`
	Enable           bool    `gorm:"column:enable"`
	Sort             *int    `gorm:"column:sort"`
	CreatedAt        int64   `gorm:"column:created_at"`
	UpdatedAt        int64   `gorm:"column:updated_at"`
}

func (Payment) TableName() string { return "v2_payment" }
