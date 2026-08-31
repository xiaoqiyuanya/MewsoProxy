package dto

type NodeReportReq struct {
	NodeType string           `json:"node_type" binding:"required"`
	NodeID   int              `json:"node_id" binding:"required,min=1"`
	Token    string           `json:"token" binding:"required"`
	Online   *bool            `json:"online"`
	Uptime   int64            `json:"uptime"`
	Load     float64          `json:"load"`
	U        int64            `json:"u"`
	D        int64            `json:"d"`
	Users    []NodeReportUser `json:"users"`
}

type NodeReportUser struct {
	UUID string `json:"uuid" binding:"required"`
	U    int64  `json:"u"`
	D    int64  `json:"d"`
}
