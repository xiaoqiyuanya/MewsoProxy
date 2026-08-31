package node

import (
	"github.com/gin-gonic/gin"

	"mewsoproxy/server/dto"
	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/response"
	nodesvc "mewsoproxy/server/service/node"
)

type Handler struct {
	nodeSvc *nodesvc.Service
}

func New(nodeSvc *nodesvc.Service) *Handler {
	return &Handler{nodeSvc: nodeSvc}
}

func (h *Handler) Report(c *gin.Context) {
	var req dto.NodeReportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	if err := h.nodeSvc.Report(c.Request.Context(), req); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}
