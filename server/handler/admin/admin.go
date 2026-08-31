package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"mewsoproxy/server/dto"
	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/response"
	"mewsoproxy/server/service/admin"
	authsvc "mewsoproxy/server/service/auth"
	usersvc "mewsoproxy/server/service/user"
)

type Handler struct {
	adminSvc *admin.Service
	authSvc  *authsvc.Service
	userSvc  *usersvc.Service
}

func New(adminSvc *admin.Service, authSvc *authsvc.Service, userSvc *usersvc.Service) *Handler {
	return &Handler{adminSvc: adminSvc, authSvc: authSvc, userSvc: userSvc}
}

func (h *Handler) Login(c *gin.Context) {
	var req dto.AdminLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	u, err := h.userSvc.Login(c.Request.Context(), dto.LoginReq{Email: req.Email, Password: req.Password})
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	if !u.IsAdmin {
		response.Fail(c, apperror.CodeNoPermission, "无管理员权限")
		return
	}
	tk, err := h.authSvc.Issue(c.Request.Context(), u)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, gin.H{"token": tk, "user": usersvc.ToDTO(u)})
}

func parsePage(c *gin.Context) (int, int) {
	page := intDefault(c.Query("page"), 1)
	size := intDefault(c.Query("page_size"), 20)
	return page, size
}

func intDefault(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
