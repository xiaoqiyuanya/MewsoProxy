package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ApiResponse{Code: 0, Message: "ok", Data: data})
}

func Fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, ApiResponse{Code: code, Message: message})
}
