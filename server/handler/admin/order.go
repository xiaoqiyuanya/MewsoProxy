package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"mewsoproxy/server/dto"
	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/response"
)

func (h *Handler) ListOrders(c *gin.Context) {
	page, size := parsePage(c)
	status := intDefault(c.Query("status"), -1)
	list, total, err := h.adminSvc.ListOrders(c.Request.Context(), int8(status), c.Query("keyword"), page, size)
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

func (h *Handler) GetOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id <= 0 {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	o, err := h.adminSvc.GetOrder(c.Request.Context(), uint(id))
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, dto.ToOrderDTO(o))
}

func (h *Handler) CancelOrder(c *gin.Context) {
	var req dto.AdminBanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	if err := h.adminSvc.CancelOrder(c.Request.Context(), req.ID); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}

func (h *Handler) MarkOrderPaid(c *gin.Context) {
	var req struct {
		ID         uint   `json:"id" binding:"required,min=1"`
		CallbackNo string `json:"callback_no"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	if err := h.adminSvc.MarkOrderPaid(c.Request.Context(), req.ID, req.CallbackNo); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}
