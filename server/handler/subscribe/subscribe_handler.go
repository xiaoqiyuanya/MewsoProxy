package subscribe

import (
	"fmt"
	"net/http"
	"strings"

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

func (h *Handler) Download(c *gin.Context) {
	token := c.Query("token")
	u, err := h.userSvc.GetByToken(c.Request.Context(), token)
	if err != nil {
		c.String(http.StatusOK, base64Error(apperror.UserMsg(err)))
		return
	}
	ua := c.GetHeader("User-Agent")
	ctype := subscribe.ContentV2ray
	if strings.Contains(strings.ToLower(ua), "clash") {
		ctype = subscribe.ContentClash
	}
	body, err := h.subSvc.BuildSubscription(c.Request.Context(), u, ctype)
	if err != nil {
		c.String(http.StatusOK, base64Error(apperror.UserMsg(err)))
		return
	}
	if ctype == subscribe.ContentClash {
		c.Header("Content-Type", "text/yaml; charset=utf-8")
	} else {
		c.Header("Content-Type", "text/plain; charset=utf-8")
	}
	c.Header("Cache-Control", "no-store")
	c.Header("Subscription-Userinfo", fmt.Sprintf("upload=%d; download=%d; total=%d; expire=%d", u.U, u.D, u.TransferEnable, u.ExpiredAt))
	c.String(http.StatusOK, body)
}

func base64Error(msg string) string {
	return subscribe.EncodeURL("Subscribe Error: " + msg)
}
