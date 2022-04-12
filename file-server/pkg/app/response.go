package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type Gin struct {
	C *gin.Context
}

func (g *Gin) Success(data interface{}) {
	g.C.JSON(http.StatusOK, Response{
		Code: 200,
		Data: data,
		Msg:  "",
	})
}

func (g *Gin) Fail() {
	g.C.JSON(http.StatusOK, Response{
		Code: 500,
		Data: nil,
		Msg:  "系统出错",
	})
}

func (g *Gin) Unauthorized() {
	g.C.JSON(http.StatusUnauthorized, Response{
		Code: 401,
		Data: nil,
		Msg:  "",
	})
}

//func PrintRes(w http.ResponseWriter, data interface{}) {
//	w.Header().Set("Content-Type", "application/json")
//	resp := Response{
//		Code: 200,
//		Data: data,
//		Msg:  "",
//	}
//	jsonStr, err := json.Marshal(&resp)
//	if err != nil {
//		return
//	}
//	w.Write(jsonStr)
//}
//
//func PrintErr(w http.ResponseWriter) {
//	w.Header().Set("Content-Type", "application/json")
//	resp := Response{
//		Code: 401,
//		Data: nil,
//		Msg:  "没有权限",
//	}
//	jsonStr, err := json.Marshal(&resp)
//	if err != nil {
//		return
//	}
//	w.Write(jsonStr)
//}
