package payment

import (
	"github.com/gin-gonic/gin"

	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/response"
	"mewsoproxy/server/service/order"
)

type Handler struct {
	orderSvc *order.Service
}

func New(orderSvc *order.Service) *Handler {
	return &Handler{orderSvc: orderSvc}
}

func (h *Handler) Notify(c *gin.Context) {
	var req struct {
		TradeNo  string `json:"trade_no" binding:"required"`
		CallbackNo string `json:"callback_no"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	if err := h.orderSvc.MarkPaidByTradeNo(c.Request.Context(), req.TradeNo, req.CallbackNo); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}
