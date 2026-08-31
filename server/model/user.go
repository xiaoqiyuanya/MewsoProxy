package model

type User struct {
	ID               uint   `gorm:"primarykey;column:id"`
	InviteUserID     *uint  `gorm:"column:invite_user_id"`
	TelegramID       *int64 `gorm:"column:telegram_id"`
	Email            string `gorm:"column:email;uniqueIndex"`
	Password         string `gorm:"column:password"`
	PasswordAlgo     *string `gorm:"column:password_algo"`
	PasswordSalt     *string `gorm:"column:password_salt"`
	Balance          int    `gorm:"column:balance"`
	Discount         *int   `gorm:"column:discount"`
	CommissionType   int8   `gorm:"column:commission_type"`
	CommissionRate   *int   `gorm:"column:commission_rate"`
	CommissionBalance int   `gorm:"column:commission_balance"`
	T                int64  `gorm:"column:t"`
	U                int64  `gorm:"column:u"`
	D                int64  `gorm:"column:d"`
	TransferEnable   int64  `gorm:"column:transfer_enable"`
	Banned           bool   `gorm:"column:banned"`
	IsAdmin          bool   `gorm:"column:is_admin"`
	LastLoginAt      *int64 `gorm:"column:last_login_at"`
	IsStaff          bool   `gorm:"column:is_staff"`
	LastLoginIP      *int64 `gorm:"column:last_login_ip"`
	UUID             string `gorm:"column:uuid"`
	GroupID          *int   `gorm:"column:group_id"`
	PlanID           *int   `gorm:"column:plan_id"`
	SpeedLimit       *int   `gorm:"column:speed_limit"`
	RemindExpire     *int8  `gorm:"column:remind_expire"`
	RemindTraffic    *int8  `gorm:"column:remind_traffic"`
	Token            string `gorm:"column:token"`
	ExpiredAt        int64  `gorm:"column:expired_at"`
	Remarks          *string `gorm:"column:remarks"`
	CreatedAt        int64  `gorm:"column:created_at"`
	UpdatedAt        int64  `gorm:"column:updated_at"`
}

func (User) TableName() string { return "v2_user" }
