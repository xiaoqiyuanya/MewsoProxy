package subscribe

import (
	"github.com/gin-gonic/gin"

	"mewsoproxy/server/middleware"
	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/response"
	"mewsoproxy/server/service/subscribe"
	usersvc "mewsoproxy/server/service/user"
)

type Handler struct {
	subSvc  *subscribe.Service
	userSvc *usersvc.Service
}

func New(subSvc *subscribe.Service, userSvc *usersvc.Service) *Handler {
	return &Handler{subSvc: subSvc, userSvc: userSvc}
}

func (h *Handler) Info(c *gin.Context) {
	uid := middleware.UserIDFrom(c)
	u, err := h.userSvc.GetByID(c.Request.Context(), uid)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, h.subSvc.Generate(c.Request.Context(), u))
}
