package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SuccessResponse(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, HttpResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func ErrorResponse(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, HttpResponse{
		Code: code,
		Msg:  msg,
	})
}
