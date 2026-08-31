package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"mewsoproxy/server/dto"
	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/response"
	"mewsoproxy/server/service/auth"
	usersvc "mewsoproxy/server/service/user"
)

type Handler struct {
	authSvc *auth.Service
	userSvc *usersvc.Service
}

func New(authSvc *auth.Service, userSvc *usersvc.Service) *Handler {
	return &Handler{authSvc: authSvc, userSvc: userSvc}
}

func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	u, err := h.userSvc.Register(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	tk, err := h.authSvc.Issue(c.Request.Context(), u)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	setRefreshCookie(c, tk.RefreshToken, 30*24*time.Hour)
	response.OK(c, gin.H{"token": tk, "user": usersvc.ToDTO(u)})
}

func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	u, err := h.userSvc.Login(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	tk, err := h.authSvc.Issue(c.Request.Context(), u)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	setRefreshCookie(c, tk.RefreshToken, 30*24*time.Hour)
	response.OK(c, gin.H{"token": tk, "user": usersvc.ToDTO(u)})
}

func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")
	if refreshToken == "" {
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}
		_ = c.ShouldBindJSON(&body)
		refreshToken = body.RefreshToken
	}
	if refreshToken == "" {
		response.Fail(c, apperror.CodeTokenExpired, "登录已过期")
		return
	}
	tk, err := h.authSvc.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	setRefreshCookie(c, tk.RefreshToken, 30*24*time.Hour)
	response.OK(c, tk)
}

func setRefreshCookie(c *gin.Context, token string, ttl time.Duration) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("refresh_token", token, int(ttl.Seconds()), "/", "", false, true)
}
