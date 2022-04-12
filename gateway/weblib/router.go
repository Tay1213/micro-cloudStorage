package weblib

import (
	"github.com/gin-gonic/gin"
	v1 "micro-cloudStorage/gateway/weblib/handlers/v1"
	"micro-cloudStorage/gateway/weblib/middleware"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors(), middleware.ErrorMiddleware())

	v := r.Group("/api/v1")
	{
		user := v.Group("/user")
		{
			user.GET("/name/:name", v1.GetUserByName)
			user.GET("/email/:email", v1.GetUserByEmail)
			user.POST("/login", v1.Login)
			user.POST("/logout", v1.Logout)
			user.POST("/reg", v1.Reg)
			user.DELETE("/delete/:id", middleware.Auth(), v1.DeleteUser)
			user.POST("/update", middleware.Auth(), v1.UpdateUser)
		}
		file := v.Group("file", middleware.Auth())
		{
			file.GET("/get/:id", v1.GetFile)
			file.POST("/getAll", v1.GetFiles)
			file.POST("/add", v1.AddFile)
			file.POST("/addDict", v1.AddDict)
			file.DELETE("/delete/:id", v1.DeleteFile)
		}
	}
	r.Use()
	return r
}
