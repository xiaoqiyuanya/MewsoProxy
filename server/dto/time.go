package dto

import "time"

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	return jsonMarshal(time.Time(t).Format(time.RFC3339))
}

func FromUnix(sec int64) Time {
	if sec == 0 {
		return Time(time.Time{})
	}
	return Time(time.Unix(sec, 0).UTC())
}

func FromTime(t time.Time) Time {
	return Time(t)
}

type PageInfo struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type PageReq struct {
	Page     int `json:"page" form:"page" binding:"omitempty,min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"omitempty,min=1,max=100"`
}
