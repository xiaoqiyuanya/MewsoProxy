package dto

import "mewsoproxy/server/model"

type AdminLoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,max=64"`
}

type AdminUserDTO struct {
	ID            uint   `json:"id"`
	Email         string `json:"email"`
	Balance       int    `json:"balance"`
	CommissionBalance int `json:"commission_balance"`
	IsAdmin       bool   `json:"is_admin"`
	Banned        bool   `json:"banned"`
	PlanID        *int   `json:"plan_id,omitempty"`
	GroupID       *int   `json:"group_id,omitempty"`
	ExpiredAt     Time   `json:"expired_at"`
	TransferEnable int64 `json:"transfer_enable"`
	UsedTraffic   int64  `json:"used_traffic"`
	CreatedAt     Time   `json:"created_at"`
}

type AdminUserUpdateReq struct {
	ID       uint `json:"id" binding:"required,min=1"`
	Balance  *int `json:"balance,omitempty"`
	GroupID  *int `json:"group_id,omitempty"`
	PlanID   *int `json:"plan_id,omitempty"`
	ExpiredAt *int64 `json:"expired_at,omitempty"`
	Banned   *bool `json:"banned,omitempty"`
}

type AdminBanReq struct {
	ID  uint `json:"id" binding:"required,min=1"`
	Ban bool `json:"ban"`
}

type AdminResetSecretReq struct {
	ID uint `json:"id" binding:"required,min=1"`
}

type AdminPlanSaveReq struct {
	ID              int    `json:"id"`
	GroupID         int    `json:"group_id" binding:"required"`
	TransferEnable  int64  `json:"transfer_enable" binding:"required"`
	Name            string `json:"name" binding:"required"`
	SpeedLimit      *int   `json:"speed_limit,omitempty"`
	Show            bool   `json:"show"`
	Sort            *int   `json:"sort,omitempty"`
	Renew           bool   `json:"renew"`
	Content         *string `json:"content,omitempty"`
	MonthPrice      *int   `json:"month_price,omitempty"`
	QuarterPrice    *int   `json:"quarter_price,omitempty"`
	HalfYearPrice   *int   `json:"half_year_price,omitempty"`
	YearPrice       *int   `json:"year_price,omitempty"`
	TwoYearPrice    *int   `json:"two_year_price,omitempty"`
	ThreeYearPrice  *int   `json:"three_year_price,omitempty"`
	OnetimePrice    *int   `json:"onetime_price,omitempty"`
	ResetPrice      *int   `json:"reset_price,omitempty"`
	ResetTrafficMethod *int8 `json:"reset_traffic_method,omitempty"`
	CapacityLimit   *int   `json:"capacity_limit,omitempty"`
}

type AdminCouponSaveReq struct {
	ID               uint   `json:"id"`
	Code             string `json:"code" binding:"required"`
	Name             string `json:"name"`
	Type             int8   `json:"type"`
	Value            int    `json:"value"`
	Show             bool   `json:"show"`
	LimitUse         *int   `json:"limit_use,omitempty"`
	LimitUseWithUser *int   `json:"limit_use_with_user,omitempty"`
	LimitPlanIDs     *string `json:"limit_plan_ids,omitempty"`
	LimitPeriod      *string `json:"limit_period,omitempty"`
	StartedAt        int64  `json:"started_at"`
	EndedAt          int64  `json:"ended_at"`
}

type AdminNoticeSaveReq struct {
	ID      uint    `json:"id"`
	Title   string  `json:"title" binding:"required"`
	Content string  `json:"content" binding:"required"`
	Show    bool    `json:"show"`
	ImgURL  *string `json:"img_url,omitempty"`
	Tags    *string `json:"tags,omitempty"`
}

type AdminPaymentSaveReq struct {
	ID          int     `json:"id"`
	UUID        string  `json:"uuid" binding:"required"`
	Payment     string  `json:"payment" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Icon        *string `json:"icon,omitempty"`
	Config      string  `json:"config"`
	NotifyDomain *string `json:"notify_domain,omitempty"`
	HandlingFeeFixed *int `json:"handling_fee_fixed,omitempty"`
	HandlingFeePercent *float64 `json:"handling_fee_percent,omitempty"`
	Enable      bool    `json:"enable"`
	Sort        *int    `json:"sort,omitempty"`
}

type AdminServerGroupSaveReq struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}

type AdminServerNodeSaveReq struct {
	ID         int    `json:"id"`
	Type       string `json:"type" binding:"required"` // trojan|vmess|shadowsocks|hysteria
	GroupID    string `json:"group_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Host       string `json:"host" binding:"required"`
	Port       string `json:"port" binding:"required"`
	ServerPort int    `json:"server_port"`
	Rate       string `json:"rate"`
	Show       bool   `json:"show"`
	Sort       *int   `json:"sort,omitempty"`
	Tags       *string `json:"tags,omitempty"`
	RouteID    *string `json:"route_id,omitempty"`
	ParentID   *int   `json:"parent_id,omitempty"`
	// 协议专属字段
	Cipher         *string `json:"cipher,omitempty"`
	Network        *string `json:"network,omitempty"`
	TLS            *int8   `json:"tls,omitempty"`
	AllowInsecure  *bool   `json:"allow_insecure,omitempty"`
	ServerName     *string `json:"server_name,omitempty"`
	UpMbps         *int    `json:"up_mbps,omitempty"`
	DownMbps       *int    `json:"down_mbps,omitempty"`
	Insecure       *bool   `json:"insecure,omitempty"`
}

type AdminConfigReq struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value"`
}

type AdminConfigDTO struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AdminSystemStatusDTO struct {
	ServerTime      int64  `json:"server_time"`
	DBStatus        string `json:"db_status"`
	RedisStatus     string `json:"redis_status"`
	OnlineUserCount int64  `json:"online_user_count"`
	UserCount       int64  `json:"user_count"`
	OrderCount      int64  `json:"order_count"`
	TodayPaidTotal  int64  `json:"today_paid_total"`
}

type AdminDropReq struct {
	ID int `json:"id" binding:"required,min=1"`
}

func ToAdminUserDTO(u *model.User) AdminUserDTO {
	return AdminUserDTO{
		ID:               u.ID,
		Email:            u.Email,
		Balance:          u.Balance,
		CommissionBalance: u.CommissionBalance,
		IsAdmin:          u.IsAdmin,
		Banned:           u.Banned,
		PlanID:           u.PlanID,
		GroupID:          u.GroupID,
		ExpiredAt:        FromUnix(u.ExpiredAt),
		TransferEnable:   u.TransferEnable,
		UsedTraffic:      u.U + u.D,
		CreatedAt:        FromUnix(u.CreatedAt),
	}
}
