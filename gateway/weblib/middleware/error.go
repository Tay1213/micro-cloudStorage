package middleware

import (
	"github.com/gin-gonic/gin"
	"micro-cloudStorage/gateway/pkg/app"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			appG := app.Gin{C: ctx}
			r := recover()
			if r != nil {
				appG.Fail()
			}
		}()
		ctx.Next()
	}

}
