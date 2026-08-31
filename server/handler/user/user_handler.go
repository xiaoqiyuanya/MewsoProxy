package user

import (
	"github.com/gin-gonic/gin"

	"mewsoproxy/server/middleware"
	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/response"
	usersvc "mewsoproxy/server/service/user"
)

type Handler struct {
	userSvc *usersvc.Service
}

func New(userSvc *usersvc.Service) *Handler {
	return &Handler{userSvc: userSvc}
}

func (h *Handler) Me(c *gin.Context) {
	uid := middleware.UserIDFrom(c)
	u, err := h.userSvc.GetByID(c.Request.Context(), uid)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, usersvc.ToDTO(u))
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	response.OK(c, nil)
}
