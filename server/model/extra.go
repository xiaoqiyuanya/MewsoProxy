package model

type Stat struct {
	ID             int    `gorm:"primarykey;column:id"`
	RecordAt       int    `gorm:"column:record_at;uniqueIndex"`
	RecordType     string `gorm:"column:record_type"`
	OrderCount     int    `gorm:"column:order_count"`
	OrderTotal     int    `gorm:"column:order_total"`
	CommissionCount int   `gorm:"column:commission_count"`
	CommissionTotal int   `gorm:"column:commission_total"`
	PaidCount      int    `gorm:"column:paid_count"`
	PaidTotal      int    `gorm:"column:paid_total"`
	RegisterCount  int    `gorm:"column:register_count"`
	InviteCount    int    `gorm:"column:invite_count"`
	TransferUsedTotal string `gorm:"column:transfer_used_total"`
	CreatedAt      int64  `gorm:"column:created_at"`
	UpdatedAt      int64  `gorm:"column:updated_at"`
}

func (Stat) TableName() string { return "v2_stat" }

type StatServer struct {
	ID         int    `gorm:"primarykey;column:id"`
	ServerID   int    `gorm:"column:server_id"`
	ServerType string `gorm:"column:server_type;type:char(11)"`
	U          int64  `gorm:"column:u"`
	D          int64  `gorm:"column:d"`
	RecordType string `gorm:"column:record_type"`
	RecordAt   int    `gorm:"column:record_at"`
	CreatedAt  int64  `gorm:"column:created_at"`
	UpdatedAt  int64  `gorm:"column:updated_at"`
}

func (StatServer) TableName() string { return "v2_stat_server" }

type StatUser struct {
	ID         int    `gorm:"primarykey;column:id"`
	UserID     int    `gorm:"column:user_id"`
	ServerRate float64 `gorm:"column:server_rate"`
	U          int64  `gorm:"column:u"`
	D          int64  `gorm:"column:d"`
	RecordType string `gorm:"column:record_type"`
	RecordAt   int    `gorm:"column:record_at"`
	CreatedAt  int64  `gorm:"column:created_at"`
	UpdatedAt  int64  `gorm:"column:updated_at"`
}

func (StatUser) TableName() string { return "v2_stat_user" }

type Ticket struct {
	ID          uint   `gorm:"primarykey;column:id"`
	UserID      uint   `gorm:"column:user_id"`
	Subject     string `gorm:"column:subject"`
	Level       int8   `gorm:"column:level"`
	Status      int8   `gorm:"column:status"`
	ReplyStatus int8   `gorm:"column:reply_status"`
	CreatedAt   int64  `gorm:"column:created_at"`
	UpdatedAt   int64  `gorm:"column:updated_at"`
}

func (Ticket) TableName() string { return "v2_ticket" }

type TicketMessage struct {
	ID        uint   `gorm:"primarykey;column:id"`
	UserID    uint   `gorm:"column:user_id"`
	TicketID  uint   `gorm:"column:ticket_id"`
	Message   string `gorm:"type:text;column:message"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

func (TicketMessage) TableName() string { return "v2_ticket_message" }

type Knowledge struct {
	ID        uint   `gorm:"primarykey;column:id"`
	Language  string `gorm:"column:language"`
	Category  string `gorm:"column:category"`
	Title     string `gorm:"column:title"`
	Body      string `gorm:"type:text;column:body"`
	Sort      *int   `gorm:"column:sort"`
	Show      bool   `gorm:"column:show"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

func (Knowledge) TableName() string { return "v2_knowledge" }

type Log struct {
	ID        uint   `gorm:"primarykey;column:id"`
	Title     string `gorm:"column:title"`
	Level     *string `gorm:"column:level"`
	Host      *string `gorm:"column:host"`
	URI       string `gorm:"column:uri"`
	Method    string `gorm:"column:method"`
	Data      *string `gorm:"type:text;column:data"`
	IP        *string `gorm:"column:ip"`
	Context   *string `gorm:"type:text;column:context"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

func (Log) TableName() string { return "v2_log" }
