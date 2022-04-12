package app

import (
	"github.com/gin-gonic/gin"
	"micro-cloudStorage/gateway/pkg/e"
	"net/http"
)

type Gin struct {
	C *gin.Context
}

type Respond struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func (g *Gin) Success(data interface{}) {
	g.C.JSON(http.StatusOK, Respond{
		Code: e.SUCCESS,
		Data: data,
		Msg:  e.GetMsg(e.SUCCESS),
	})
}

func (g *Gin) Fail() {
	g.C.JSON(http.StatusOK, Respond{
		Code: 500,
		Data: nil,
		Msg:  "系统出错",
	})
}

func (g *Gin) FailMsg(code int, msg string) {
	if msg == "" {
		msg = e.GetMsg(code)
	}
	g.C.JSON(http.StatusOK, Respond{
		Code: code,
		Data: nil,
		Msg:  msg,
	})
}

func (g *Gin) Unauthorized(code int) {
	g.C.JSON(http.StatusUnauthorized, Respond{
		Code: code,
		Data: nil,
		Msg:  e.GetMsg(code),
	})
}
