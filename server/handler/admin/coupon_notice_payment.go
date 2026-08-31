package admin

import (
	"github.com/gin-gonic/gin"

	"mewsoproxy/server/dto"
	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/response"
)

func (h *Handler) ListCoupons(c *gin.Context) {
	list, err := h.adminSvc.ListCoupons(c.Request.Context())
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, list)
}

func (h *Handler) SaveCoupon(c *gin.Context) {
	var req dto.AdminCouponSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	id, err := h.adminSvc.SaveCoupon(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *Handler) DropCoupon(c *gin.Context) {
	var req dto.AdminDropReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	if err := h.adminSvc.DropCoupon(c.Request.Context(), uint(req.ID)); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}

func (h *Handler) ToggleCouponShow(c *gin.Context) {
	var req struct {
		ID   uint   `json:"id" binding:"required,min=1"`
		Show bool   `json:"show"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	if err := h.adminSvc.ToggleCouponShow(c.Request.Context(), req.ID, req.Show); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}

func (h *Handler) ListNotices(c *gin.Context) {
	list, err := h.adminSvc.ListNotices(c.Request.Context())
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, list)
}

func (h *Handler) SaveNotice(c *gin.Context) {
	var req dto.AdminNoticeSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	id, err := h.adminSvc.SaveNotice(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *Handler) DropNotice(c *gin.Context) {
	var req dto.AdminDropReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	if err := h.adminSvc.DropNotice(c.Request.Context(), uint(req.ID)); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}

func (h *Handler) ToggleNoticeShow(c *gin.Context) {
	var req struct {
		ID   uint `json:"id" binding:"required,min=1"`
		Show bool `json:"show"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	if err := h.adminSvc.ToggleNoticeShow(c.Request.Context(), req.ID, req.Show); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}

func (h *Handler) ListPayments(c *gin.Context) {
	list, err := h.adminSvc.ListPayments(c.Request.Context())
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, list)
}

func (h *Handler) SavePayment(c *gin.Context) {
	var req dto.AdminPaymentSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	id, err := h.adminSvc.SavePayment(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *Handler) DropPayment(c *gin.Context) {
	var req dto.AdminDropReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	if err := h.adminSvc.DropPayment(c.Request.Context(), req.ID); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}

func (h *Handler) TogglePaymentShow(c *gin.Context) {
	var req struct {
		ID     int  `json:"id" binding:"required,min=1"`
		Enable bool `json:"enable"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	if err := h.adminSvc.TogglePaymentShow(c.Request.Context(), req.ID, req.Enable); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}
