package admin

import (
	"github.com/gin-gonic/gin"

	"mewsoproxy/server/dto"
	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/response"
	"mewsoproxy/server/service/plan"
)

func (h *Handler) ListPlans(c *gin.Context) {
	list, err := h.adminSvc.ListPlans(c.Request.Context())
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	out := make([]dto.PlanDetailDTO, 0, len(list))
	for i := range list {
		out = append(out, plan.ToDTO(&list[i]))
	}
	response.OK(c, out)
}

func (h *Handler) SavePlan(c *gin.Context) {
	var req dto.AdminPlanSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	id, err := h.adminSvc.SavePlan(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *Handler) DropPlan(c *gin.Context) {
	var req dto.AdminDropReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	if err := h.adminSvc.DropPlan(c.Request.Context(), req.ID); err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, nil)
}
