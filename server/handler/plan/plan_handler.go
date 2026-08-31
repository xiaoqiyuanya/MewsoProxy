package plan

import (
	"github.com/gin-gonic/gin"

	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/response"
	"mewsoproxy/server/service/plan"
)

type Handler struct {
	planSvc *plan.Service
}

func New(planSvc *plan.Service) *Handler {
	return &Handler{planSvc: planSvc}
}

func (h *Handler) List(c *gin.Context) {
	list, err := h.planSvc.List(c.Request.Context())
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, list)
}

func (h *Handler) Get(c *gin.Context) {
	var req struct {
		ID int `json:"id" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	p, err := h.planSvc.Get(c.Request.Context(), req.ID)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, plan.ToDTO(p))
}
