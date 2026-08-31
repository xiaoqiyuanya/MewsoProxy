package admin

import (
	"github.com/gin-gonic/gin"

	"mewsoproxy/server/dto"
	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/response"
)

func (h *Handler) ListUsers(c *gin.Context) {
	page, size := parsePage(c)
	list, total, err := h.adminSvc.ListUsers(c.Request.Context(), c.Query("keyword"), page, size)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	out := make([]dto.AdminUserDTO, 0, len(list))
	for i := range list {
		out = append(out, dto.ToAdminUserDTO(&list[i]))
	}
	response.OK(c, gin.H{"list": out, "total": total})
}

func (h *Handler) GetUser(c *gin.Context) {
	var req dto.AdminResetSecretReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	u, err := h.adminSvc.GetUser(c.Request.Context(), req.ID)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, dto.ToAdminUserDTO(u))
}

func (h *Handler) UpdateUser(c *gin.Context) {
	var req dto.AdminUserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	if err := h.adminSvc.UpdateUser(c.Request.Context(), req); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}

func (h *Handler) BanUser(c *gin.Context) {
	var req dto.AdminBanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	if err := h.adminSvc.BanUser(c.Request.Context(), req.ID, req.Ban); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}

func (h *Handler) ResetSecret(c *gin.Context) {
	var req dto.AdminResetSecretReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	if err := h.adminSvc.ResetSecret(c.Request.Context(), req.ID); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}
