package admin

import (
	"github.com/gin-gonic/gin"

	"mewsoproxy/server/dto"
	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/response"
)

func (h *Handler) ListGroups(c *gin.Context) {
	list, err := h.adminSvc.ListGroups(c.Request.Context())
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, list)
}

func (h *Handler) SaveGroup(c *gin.Context) {
	var req dto.AdminServerGroupSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	id, err := h.adminSvc.SaveGroup(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *Handler) DropGroup(c *gin.Context) {
	var req dto.AdminDropReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	if err := h.adminSvc.DropGroup(c.Request.Context(), req.ID); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}

func (h *Handler) ListNodes(c *gin.Context) {
	nodeType := c.Query("type")
	list, err := h.adminSvc.ListNodes(c.Request.Context(), nodeType)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, list)
}

func (h *Handler) SaveNode(c *gin.Context) {
	var req dto.AdminServerNodeSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	id, err := h.adminSvc.SaveNode(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *Handler) DropNode(c *gin.Context) {
	var req struct {
		ID   int    `json:"id" binding:"required,min=1"`
		Type string `json:"type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	if err := h.adminSvc.DropNode(c.Request.Context(), req.Type, req.ID); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}
