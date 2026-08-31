package dto

type UserDTO struct {
	ID            uint   `json:"id"`
	Email         string `json:"email"`
	Balance       int    `json:"balance"`
	CommissionBalance int `json:"commission_balance"`
	IsAdmin       bool   `json:"is_admin"`
	IsStaff       bool   `json:"is_staff"`
	Banned        bool   `json:"banned"`
	PlanID        *int   `json:"plan_id,omitempty"`
	GroupID       *int   `json:"group_id,omitempty"`
	ExpiredAt     Time   `json:"expired_at"`
	Token         string `json:"token"`
	UUID          string `json:"uuid"`
	TransferEnable int64 `json:"transfer_enable"`
	UsedTraffic   int64  `json:"used_traffic"`
	CreatedAt     Time   `json:"created_at"`
}

type UserInfoReq struct {
}

type UpdateProfileReq struct {
	Name string `json:"email" binding:"omitempty,email"`
}
