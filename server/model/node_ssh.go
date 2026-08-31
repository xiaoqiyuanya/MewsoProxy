package model

const (
	NodeInstallStatusNone     int8 = 0
	NodeInstallStatusRunning  int8 = 1
	NodeInstallStatusSuccess  int8 = 2
	NodeInstallStatusFailed   int8 = 3
)

type ServerNodeSSH struct {
	ID            uint    `gorm:"primarykey;column:id"`
	NodeType      string  `gorm:"column:node_type;size:32"`
	NodeID        int     `gorm:"column:node_id"`
	SSHHost       string  `gorm:"column:ssh_host;size:255"`
	SSHPort       int     `gorm:"column:ssh_port"`
	SSHUser       string  `gorm:"column:ssh_user;size:64"`
	SSHPassword   *string `gorm:"column:ssh_password;type:text"`
	SSHPrivateKey *string `gorm:"column:ssh_private_key;type:text"`
	InstallStatus int8    `gorm:"column:install_status"`
	LastLog       *string `gorm:"column:last_log;type:text"`
	CreatedAt     int64   `gorm:"column:created_at"`
	UpdatedAt     int64   `gorm:"column:updated_at"`
}

func (ServerNodeSSH) TableName() string { return "v2_server_node_ssh" }
