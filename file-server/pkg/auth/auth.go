package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"micro-cloudStorage/fileServer/pkg/app"
	"micro-cloudStorage/fileServer/pkg/constant"
	"micro-cloudStorage/fileServer/pkg/gredis"
	"strconv"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		appG := app.Gin{C: context}
		auth := context.Request.Header.Get("Authorization")
		if auth == "" {
			appG.Fail()
			context.Abort()
			return
		}
		//切割字符串，把Bearer前缀去掉
		auth = strings.Fields(auth)[1]
		//解析jwt
		claims, err := parseToken(auth)
		if err != nil {
			appG.Unauthorized()
			context.Abort()
			return
		}
		//查询redis，token是否存在
		res := gredis.Exists(constant.TOKEN_PRIFIX + ":" + claims.Id + ":" + auth)
		if res {
			fmt.Println("token正确")
			id, _ := strconv.Atoi(claims.Id)
			context.Set("id", id)
		} else {
			appG.Unauthorized()
			context.Abort()
			return
		}
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
