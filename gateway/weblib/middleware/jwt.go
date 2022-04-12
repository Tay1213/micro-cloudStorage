package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"micro-cloudStorage/gateway/pkg/app"
	"micro-cloudStorage/gateway/pkg/constant"
	"micro-cloudStorage/gateway/pkg/e"
	"micro-cloudStorage/gateway/pkg/gredis"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		appG := app.Gin{C: context}

		//获取header上的Authentication参数
		auth := context.Request.Header.Get("Authorization")
		if auth == "" {
			context.Abort()
			appG.Unauthorized(e.TOKEN_INVALID)
		}
		//切割字符串，把Bearer前缀去掉
		auth = strings.Fields(auth)[1]
		//解析jwt
		claims, err := parseToken(auth)
		if err != nil {
			context.Abort()
			appG.Unauthorized(e.TOKEN_INVALID)
		}
		//查询redis，token是否存在
		res := gredis.Exists(constant.TOKEN_PRIFIX + ":" + claims.Id + ":" + auth)
		if res {
			fmt.Println("token正确")
			gredis.Expire(constant.TOKEN_PRIFIX+":"+claims.Id+":"+auth, constant.EXPIRE_TIME)
		} else {
			context.Abort()
			appG.Unauthorized(e.TOKEN_INVALID)
		}
		context.Set("userId", claims.Id)
		context.Next()
	}
}

func parseToken(token string) (*jwt.StandardClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(constant.JWT_SECRET), nil
	})
	if err == nil && jwtToken != nil {
		if claim, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
			return claim, nil
		}
	}
	return nil, err
}
