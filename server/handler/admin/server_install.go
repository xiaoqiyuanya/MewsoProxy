package admin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"mewsoproxy/server/dto"
	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/response"
)

func (h *Handler) NodeInstall(c *gin.Context) {
	var req dto.AdminNodeInstallReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.CodeParamFormat, "参数不合法")
		return
	}
	taskID, err := h.adminSvc.InstallNode(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	response.OK(c, dto.AdminNodeInstallResp{TaskID: taskID})
}

func (h *Handler) NodeInstallLog(c *gin.Context) {
	taskID := c.Query("task_id")
	ch, unsub, err := h.adminSvc.SubscribeLog(taskID)
	if err != nil {
		response.Fail(c, apperror.Code(err), apperror.UserMsg(err))
		return
	}
	defer unsub()
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.String(http.StatusOK, "SSE 不可用")
		return
	}
	flusher.Flush()
	for {
		select {
		case <-c.Request.Context().Done():
			return
		case line, okCh := <-ch:
			if !okCh {
				return
			}
			fmt.Fprintf(c.Writer, "data: %s\n\n", line)
			flusher.Flush()
			if line == "##INSTALL_DONE##" || strings.HasPrefix(line, "##INSTALL_FAILED") {
				return
			}
		}
	}
}
