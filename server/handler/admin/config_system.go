package admin

import (
	"github.com/gin-gonic/gin"

	"mewsoproxy/server/dto"
	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/response"
)

func (h *Handler) GetConfig(c *gin.Context) {
	list, err := h.adminSvc.GetConfig(c.Request.Context())
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, list)
}

func (h *Handler) SaveConfig(c *gin.Context) {
	var req dto.AdminConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	if err := h.adminSvc.SaveConfig(c.Request.Context(), req); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}

func (h *Handler) SystemStatus(c *gin.Context) {
	res, err := h.adminSvc.SystemStatus(c.Request.Context())
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, res)
}
