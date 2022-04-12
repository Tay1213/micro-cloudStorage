package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"micro-cloudStorage/gateway/pkg/app"
	"micro-cloudStorage/gateway/service"
	"strconv"
)

func Login(c *gin.Context) {
	appG := app.Gin{C: c}
	var userReq service.UserRequest
	PanicIfUserError(c.Bind(&userReq))
	userResp, err := UserServices.UserLogin(context.Background(), &userReq)
	PanicIfUserError(err)
	appG.Success(userResp)
}

func Logout(c *gin.Context) {
	appG := app.Gin{C: c}
	var userReq service.UserRequest
	PanicIfUserError(c.Bind(&userReq))
	userResp, err := UserServices.UserLogout(context.Background(), &userReq)
	PanicIfUserError(err)
	appG.Success(userResp)
}

func GetUserByName(c *gin.Context) {
	appG := app.Gin{C: c}
	name := appG.C.Param("name")
	var userReq service.UserRequest
	userReq.UserName = name
	userResp, err := UserServices.GetUserByName(context.Background(), &userReq)
	PanicIfUserError(err)
	appG.Success(map[string]interface{}{
		"client_random_value": userResp.UserDetail.ClientRandomValue,
	})

}

func GetUserByEmail(c *gin.Context) {
	appG := app.Gin{C: c}
	email := appG.C.Param("email")
	var userReq service.UserRequest
	userReq.Email = email
	userResp, err := UserServices.GetUserByEmail(context.Background(), &userReq)
	PanicIfUserError(err)
	appG.Success(map[string]interface{}{
		"client_random_value": userResp.UserDetail.ClientRandomValue,
	})
}

func Reg(c *gin.Context) {
	appG := app.Gin{C: c}
	var userReq service.UserRequest
	PanicIfUserError(c.Bind(&userReq))
	userResp, err := UserServices.UserRegister(context.Background(), &userReq)
	PanicIfUserError(err)
	appG.Success(userResp)

}

func UpdateUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var userReq service.UserRequest
	PanicIfFileError(c.Bind(&userReq))
	userResp, err := UserServices.UpdateUser(context.Background(), &userReq)
	PanicIfUserError(err)
	appG.Success(userResp)
}

func DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var userReq service.UserRequest
	id, err := strconv.Atoi(c.Param("id"))
	PanicIfUserError(err)
	userReq.ID = int32(id)
	userResp, err := UserServices.DeleteUser(context.Background(), &userReq)
	PanicIfUserError(err)
	appG.Success(userResp)
}
