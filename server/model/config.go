package model

type Config struct {
	ID     uint   `gorm:"primarykey;column:id"`
	Key       string `gorm:"column:key;type:varchar(255);uniqueIndex"`
	Value  string `gorm:"type:text;column:value"`
	CreatedAt int64 `gorm:"column:created_at"`
	UpdatedAt int64 `gorm:"column:updated_at"`
}

func (Config) TableName() string { return "v2_config" }

const (
	ConfigKeySiteName       = "app_name"
	ConfigKeySecurePath     = "secure_path"
	ConfigKeyFrontendAdminPath = "frontend_admin_path"
	ConfigKeyRegisterEnabled = "register_enabled"
	ConfigKeySubscribeURL   = "subscribe_url"
)
