package model

type ServerGroup struct {
	ID        int    `json:"id" gorm:"primarykey;column:id"`
	Name      string `json:"name" gorm:"column:name"`
	CreatedAt int64  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt int64  `json:"updated_at" gorm:"column:updated_at"`
}

func (ServerGroup) TableName() string { return "v2_server_group" }

type ServerRoute struct {
	ID          int    `gorm:"primarykey;column:id"`
	Remarks     string `gorm:"column:remarks"`
	Match       string `gorm:"type:text;column:match"`
	Action      string `gorm:"column:action"`
	ActionValue *string `gorm:"column:action_value"`
	CreatedAt   int64  `gorm:"column:created_at"`
	UpdatedAt   int64  `gorm:"column:updated_at"`
}

func (ServerRoute) TableName() string { return "v2_server_route" }

type ServerTrojan struct {
	ID            int     `json:"id" gorm:"primarykey;column:id"`
	GroupID       string  `json:"group_id" gorm:"column:group_id"`
	RouteID       *string `json:"route_id" gorm:"column:route_id"`
	ParentID      *int    `json:"parent_id" gorm:"column:parent_id"`
	Tags          *string `json:"tags" gorm:"column:tags"`
	Name          string  `json:"name" gorm:"column:name"`
	Rate          string  `json:"rate" gorm:"column:rate"`
	Host          string  `json:"host" gorm:"column:host"`
	Port          string  `json:"port" gorm:"column:port"`
	ServerPort    int     `json:"server_port" gorm:"column:server_port"`
	AllowInsecure bool    `json:"allow_insecure" gorm:"column:allow_insecure"`
	ServerName    *string `json:"server_name" gorm:"column:server_name"`
	Show          bool    `json:"show" gorm:"column:show"`
	Sort          *int    `json:"sort" gorm:"column:sort"`
	CreatedAt     int64   `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     int64   `json:"updated_at" gorm:"column:updated_at"`
}

func (ServerTrojan) TableName() string { return "v2_server_trojan" }

type ServerVmess struct {
	ID              int     `json:"id" gorm:"primarykey;column:id"`
	GroupID         string  `json:"group_id" gorm:"column:group_id"`
	RouteID         *string `json:"route_id" gorm:"column:route_id"`
	Name            string  `json:"name" gorm:"column:name"`
	ParentID        *int    `json:"parent_id" gorm:"column:parent_id"`
	Host            string  `json:"host" gorm:"column:host"`
	Port            string  `json:"port" gorm:"column:port"`
	ServerPort      int     `json:"server_port" gorm:"column:server_port"`
	TLS             int8    `json:"tls" gorm:"column:tls"`
	Tags            *string `json:"tags" gorm:"column:tags"`
	Rate            string  `json:"rate" gorm:"column:rate"`
	Network         string  `json:"network" gorm:"column:network"`
	Rules           *string `json:"rules" gorm:"type:text;column:rules"`
	NetworkSettings *string `json:"network_settings" gorm:"type:text;column:network_settings"`
	TLSSettings     *string `json:"tls_settings" gorm:"type:text;column:tls_settings"`
	RuleSettings    *string `json:"rule_settings" gorm:"type:text;column:rule_settings"`
	DNSSettings     *string `json:"dns_settings" gorm:"type:text;column:dns_settings"`
	Show            bool    `json:"show" gorm:"column:show"`
	Sort            *int    `json:"sort" gorm:"column:sort"`
	CreatedAt       int64   `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       int64   `json:"updated_at" gorm:"column:updated_at"`
}

func (ServerVmess) TableName() string { return "v2_server_vmess" }

type ServerShadowsocks struct {
	ID           int     `json:"id" gorm:"primarykey;column:id"`
	GroupID      string  `json:"group_id" gorm:"column:group_id"`
	RouteID      *string `json:"route_id" gorm:"column:route_id"`
	ParentID     *int    `json:"parent_id" gorm:"column:parent_id"`
	Tags         *string `json:"tags" gorm:"column:tags"`
	Name         string  `json:"name" gorm:"column:name"`
	Rate         string  `json:"rate" gorm:"column:rate"`
	Host         string  `json:"host" gorm:"column:host"`
	Port         string  `json:"port" gorm:"column:port"`
	ServerPort   int     `json:"server_port" gorm:"column:server_port"`
	Cipher       string  `json:"cipher" gorm:"column:cipher"`
	Obfs         *string `json:"obfs" gorm:"column:obfs"`
	ObfsSettings *string `json:"obfs_settings" gorm:"column:obfs_settings"`
	Show         bool    `json:"show" gorm:"column:show"`
	Sort         *int    `json:"sort" gorm:"column:sort"`
	CreatedAt    int64   `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    int64   `json:"updated_at" gorm:"column:updated_at"`
}

func (ServerShadowsocks) TableName() string { return "v2_server_shadowsocks" }

type ServerHysteria struct {
	ID         int     `json:"id" gorm:"primarykey;column:id"`
	GroupID    string  `json:"group_id" gorm:"column:group_id"`
	RouteID    *string `json:"route_id" gorm:"column:route_id"`
	ParentID   *int    `json:"parent_id" gorm:"column:parent_id"`
	Name       string  `json:"name" gorm:"column:name"`
	Host       string  `json:"host" gorm:"column:host"`
	Port       string  `json:"port" gorm:"column:port"`
	ServerPort int     `json:"server_port" gorm:"column:server_port"`
	Tags       *string `json:"tags" gorm:"column:tags"`
	Rate       string  `json:"rate" gorm:"column:rate"`
	Show       bool    `json:"show" gorm:"column:show"`
	Sort       *int    `json:"sort" gorm:"column:sort"`
	UpMbps     int     `json:"up_mbps" gorm:"column:up_mbps"`
	DownMbps   int     `json:"down_mbps" gorm:"column:down_mbps"`
	ServerName *string `json:"server_name" gorm:"column:server_name"`
	Insecure   bool    `json:"insecure" gorm:"column:insecure"`
	CreatedAt  int64   `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  int64   `json:"updated_at" gorm:"column:updated_at"`
}

func (ServerHysteria) TableName() string { return "v2_server_hysteria" }
