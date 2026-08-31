package payment

import (
	"github.com/gin-gonic/gin"

	"mewsoproxy/server/middleware"
	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/response"
	"mewsoproxy/server/service/payment"
)

type Handler struct {
	paymentSvc *payment.Service
}

func New(paymentSvc *payment.Service) *Handler {
	return &Handler{paymentSvc: paymentSvc}
}

func (h *Handler) Create(c *gin.Context) {
	var req struct {
		OrderID   uint `json:"order_id" binding:"required,min=1"`
		PaymentID int  `json:"payment_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	uid := middleware.UserIDFrom(c)
	baseURL := baseURLFromRequest(c)
	res, err := h.paymentSvc.Create(c.Request.Context(), uid, req.OrderID, req.PaymentID, baseURL)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, res)
}

func (h *Handler) Channels(c *gin.Context) {
	list, err := h.paymentSvc.Channels(c.Request.Context())
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, list)
}

func (h *Handler) Notify(c *gin.Context) {
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	if len(params) == 0 {
		var body map[string]string
		if err := c.ShouldBindJSON(&body); err != nil {
			response.Fail(c, apperror.CodeParamFormat, "参数不合法")
			return
		}
		params = body
	}
	if err := h.paymentSvc.Notify(c.Request.Context(), params); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}

func baseURLFromRequest(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
}
