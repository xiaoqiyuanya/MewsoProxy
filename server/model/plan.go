package model

type Plan struct {
	ID                  int    `gorm:"primarykey;column:id"`
	GroupID             int    `gorm:"column:group_id"`
	TransferEnable      int64  `gorm:"column:transfer_enable"`
	Name                string `gorm:"column:name"`
	SpeedLimit          *int   `gorm:"column:speed_limit"`
	Show                bool   `gorm:"column:show"`
	Sort                *int   `gorm:"column:sort"`
	Renew               bool   `gorm:"column:renew"`
	Content             *string `gorm:"column:content"`
	MonthPrice          *int   `gorm:"column:month_price"`
	QuarterPrice        *int   `gorm:"column:quarter_price"`
	HalfYearPrice       *int   `gorm:"column:half_year_price"`
	YearPrice           *int   `gorm:"column:year_price"`
	TwoYearPrice        *int   `gorm:"column:two_year_price"`
	ThreeYearPrice      *int   `gorm:"column:three_year_price"`
	OnetimePrice        *int   `gorm:"column:onetime_price"`
	ResetPrice          *int   `gorm:"column:reset_price"`
	ResetTrafficMethod  *int8  `gorm:"column:reset_traffic_method"`
	CapacityLimit       *int   `gorm:"column:capacity_limit"`
	CreatedAt           int64  `gorm:"column:created_at"`
	UpdatedAt           int64  `gorm:"column:updated_at"`
}

func (Plan) TableName() string { return "v2_plan" }
