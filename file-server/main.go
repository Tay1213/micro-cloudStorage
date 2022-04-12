package main

import (
	"github.com/gin-gonic/gin"
	"micro-cloudStorage/fileServer/conf"
	"micro-cloudStorage/fileServer/pkg/auth"
	"micro-cloudStorage/fileServer/pkg/gredis"
	"micro-cloudStorage/fileServer/pkg/router"
)

func main() {
	conf.Setup()
	_ = gredis.Setup()

	g := gin.Default()
	g.Use(auth.Auth())
	g.POST("/file/get", router.GetFile)

	g.Run(":9090")

}
