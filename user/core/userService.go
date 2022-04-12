package core

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"micro-cloudStorage/user/model"
	"micro-cloudStorage/user/pkg/constant"
	"micro-cloudStorage/user/pkg/gredis"
	"micro-cloudStorage/user/pkg/utils"
	"micro-cloudStorage/user/service"
	"strconv"
	"time"
)

type UserService struct {
}

func (*UserService) UserRegister(ctx context.Context, req *service.UserRequest, resp *service.UserDetailResponse) error {
	_, err := model.GetUserByName(req.UserName)
	if err == nil {
		return errors.New("用户名已被注册！")
	}
	_, err = model.GetUserByEmail(req.Email)
	if err == nil {
		return errors.New("邮箱已被注册！")
	}
	var fileReq = service.FileRequest{}
	fileResp, err := FileService.CreateNewFileRoot(ctx, &fileReq)
	if err != nil {
		return errors.New("调用文件服务失败")
	}
	_, err = model.CreateUser(req, int(fileResp.RootDictID))
	if err != nil {
		return errors.New("注册用户失败！")
	}
	return nil
}
func (*UserService) UserLogin(ctx context.Context, req *service.UserRequest, resp *service.UserDetailResponse) error {
	//验证
	user, err := model.Validate(req)
	if err != nil || user.Id == 0 {
		return errors.New("验证失败！")
	}
	//判断是否30秒内登录次数大于5
	flag, _ := gredis.Get(constant.LOGIN_PRIFIX + ":" + strconv.Itoa(user.Id))
	if flag != nil {
		times := flag[0] - 48
		if times >= 5 {
			return errors.New("登录次数过多！")
		} else {
			err = gredis.Set(constant.LOGIN_PRIFIX+":"+strconv.Itoa(user.Id), times+1, 30)
			if err != nil {
				return errors.New("设置redis出错！")
			}
		}
	} else {
		err = gredis.Set(constant.LOGIN_PRIFIX+":"+strconv.Itoa(user.Id), 1, 30)
		if err != nil {
			return errors.New("设置redis出错！")
		}
	}

	//获取token
	token, err := getJwt(user)
	if err != nil {
		return errors.New("获取token失败！")
	}
	resp.Token = "Bearer " + token
	resp.UserDetail = &service.UserModel{
		ID:                 uint32(user.Id),
		EncryptedMasterKey: user.EncryptedMasterKey,
		RootDictId:         int64(user.RootDictId),
	}

	var fileReq = &service.FileRequest{
		UserID:   int32(user.Id),
		ID:       int32(user.RootDictId),
		PageNum:  1,
		PageSize: 10,
	}
	//调用fileService，获取文件信息
	fileResp, err := FileService.GetFilesByParentDictId(context.Background(), fileReq)
	if err != nil {
		return err
	}
	resp.FileResponse = &service.FileResponse{
		Files:     fileResp.Files,
		TotalFile: fileResp.TotalFile,
	}
	//增加登录次数
	user, err = model.Login(req)
	if err != nil {
		return errors.New("登录失败！")
	}
	return nil

}
func (*UserService) UserLogout(ctx context.Context, req *service.UserRequest, resp *service.UserDetailResponse) error {
	return nil
}
func (*UserService) GetUserByName(ctx context.Context, req *service.UserRequest, resp *service.UserDetailResponse) error {
	user, err := model.GetUserByName(req.UserName)
	if err != nil {
		resp.UserDetail = &service.UserModel{
			ClientRandomValue: utils.Random16(32),
		}
		return nil
	}
	resp.UserDetail = &service.UserModel{
		ClientRandomValue: user.ClientRandomValue,
	}
	return nil
}
func (*UserService) GetUserByEmail(ctx context.Context, req *service.UserRequest, resp *service.UserDetailResponse) error {
	user, err := model.GetUserByEmail(req.Email)
	if err != nil {
		resp.UserDetail = &service.UserModel{
			ClientRandomValue: utils.Random16(32),
		}
		return nil
	}
	resp.UserDetail = &service.UserModel{
		ClientRandomValue: user.ClientRandomValue,
	}
	return nil

}

func (*UserService) UpdateUser(ctx context.Context, req *service.UserRequest, resp *service.UserDetailResponse) error {
	return model.UpdateUser(req)
}
func (*UserService) DeleteUser(ctx context.Context, req *service.UserRequest, resp *service.UserDetailResponse) error {
	return model.DeleteUser(int(req.ID))

}

func getJwt(user *model.User) (string, error) {
	id := strconv.Itoa(user.Id)
	claims := jwt.StandardClaims{
		Audience:  user.Username,                                         // 受众
		ExpiresAt: time.Now().Unix() + (int64(constant.JWT_EXPIRE_TIME)), // 失效时间
		Id:        id,                                                    // 编号
		IssuedAt:  time.Now().Unix(),                                     // 签发时间
		Issuer:    "cloud_storage",                                       // 签发人
		NotBefore: time.Now().Unix(),                                     // 生效时间
		Subject:   "login",                                               // 主题
	}
	var jwtSecret = []byte(constant.JWT_SECRET)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	err = gredis.Set(constant.TOKEN_PRIFIX+":"+id+":"+token, "token", constant.EXPIRE_TIME)
	if err != nil {
		return "", err
	}
	return token, nil
}
