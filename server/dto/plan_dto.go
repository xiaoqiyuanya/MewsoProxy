package dto

type PlanDetailDTO struct {
	ID             int   `json:"id"`
	Name           string `json:"name"`
	GroupID        int   `json:"group_id"`
	TransferEnable int64 `json:"transfer_enable"`
	SpeedLimit     *int  `json:"speed_limit,omitempty"`
	Show           bool  `json:"show"`
	Sort           *int  `json:"sort,omitempty"`
	Renew          bool  `json:"renew"`
	Content        *string `json:"content,omitempty"`
	MonthPrice     *int  `json:"month_price,omitempty"`
	QuarterPrice   *int  `json:"quarter_price,omitempty"`
	HalfYearPrice  *int  `json:"half_year_price,omitempty"`
	YearPrice      *int  `json:"year_price,omitempty"`
	TwoYearPrice   *int  `json:"two_year_price,omitempty"`
	ThreeYearPrice *int  `json:"three_year_price,omitempty"`
	OnetimePrice   *int  `json:"onetime_price,omitempty"`
}

type PlanListReq struct {
	PageReq
}
