package order

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"mewsoproxy/server/dto"
	"mewsoproxy/server/middleware"
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

func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	uid := middleware.UserIDFrom(c)
	o, err := h.orderSvc.Create(c.Request.Context(), uid, req)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, dto.ToOrderDTO(o))
}

func (h *Handler) List(c *gin.Context) {
	uid := middleware.UserIDFrom(c)
	page := parseInt(c.Query("page"), 1)
	size := parseInt(c.Query("page_size"), 20)
	list, total, err := h.orderSvc.List(c.Request.Context(), uid, page, size)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	out := make([]dto.OrderDTO, 0, len(list))
	for i := range list {
		out = append(out, dto.ToOrderDTO(&list[i]))
	}
	response.OK(c, gin.H{"list": out, "total": total})
}

func (h *Handler) Detail(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	uid := middleware.UserIDFrom(c)
	o, err := h.orderSvc.Detail(c.Request.Context(), uid, req.ID)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, dto.ToOrderDTO(o))
}

func (h *Handler) Cancel(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	uid := middleware.UserIDFrom(c)
	if err := h.orderSvc.Cancel(c.Request.Context(), uid, req.ID); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}

func parseInt(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
