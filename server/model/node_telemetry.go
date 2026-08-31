package model

type ServerNodeTelemetry struct {
	ID           int     `json:"id" gorm:"primarykey;column:id"`
	NodeType     string  `json:"node_type" gorm:"column:node_type;type:char(11);uniqueIndex:idx_node_identity"`
	NodeID       int     `json:"node_id" gorm:"column:node_id;uniqueIndex:idx_node_identity"`
	Online       bool    `json:"online" gorm:"column:online"`
	LastOnlineAt int64   `json:"last_online_at" gorm:"column:last_online_at"`
	Uptime       int64   `json:"uptime" gorm:"column:uptime"`
	Load         float64 `json:"load" gorm:"column:load"`
	CreatedAt    int64   `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    int64   `json:"updated_at" gorm:"column:updated_at"`
}

func (ServerNodeTelemetry) TableName() string { return "v2_server_node_telemetry" }
