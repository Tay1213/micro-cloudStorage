package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"micro-cloudStorage/user/service"
	"time"
)

type User struct {
	Id                      int
	Username                string
	Email                   string
	ClientRandomValue       string
	EncryptedMasterKey      string
	HashedAuthenticationKey string
	Regdate                 time.Time
	Logins                  int
	RootDictId              int
}

func GetUserByName(username string) (*User, error) {
	var user = &User{}
	err := DB.Where("username = ?", username).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}

func GetUserByEmail(email string) (*User, error) {
	var user = &User{}
	err := DB.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}

func Validate(req *service.UserRequest) (*User, error) {
	var user = &User{}
	if req.UserName != "" {
		err := DB.Where("username = ? and hashed_authentication_key = ?", req.UserName, req.HashedAuthenticationKey).
			Find(&user).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := DB.Where("email = ? and hashed_authentication_key = ?", req.Email, req.HashedAuthenticationKey).
			Find(&user).Error
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func Login(req *service.UserRequest) (*User, error) {
	var user = &User{}
	var err error
	err = DB.Model(&user).
		Where("hashed_authentication_key = ?", req.HashedAuthenticationKey).
		Update("logins", gorm.Expr("logins+1")).Error
	if err != nil {
		return nil, err
	}
	return user, err
}

func CreateUser(req *service.UserRequest, rootId int) (*User, error) {
	var user = &User{
		Username:                req.UserName,
		Email:                   req.Email,
		ClientRandomValue:       req.ClientRandomValue,
		EncryptedMasterKey:      req.EncryptedMasterKey,
		HashedAuthenticationKey: req.HashedAuthenticationKey,
		Regdate:                 time.Now(),
		Logins:                  0,
		RootDictId:              rootId,
	}

	err := DB.Model(&User{}).Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUser(req *service.UserRequest) error {
	var user = User{}
	var err error
	err = DB.Where("id = ?", req.ID).Take(&user).Error
	if err != nil {
		return err
	}

	if req.Email != "" {
		err = DB.Where("email = ?", req.Email).Find(&User{}).Error
		if err == nil {
			return errors.New("邮箱已存在")
		}
		user.Email = req.Email
	}

	if req.UserName != "" {
		err = DB.Where("username = ?", req.UserName).Find(&User{}).Error
		if err == nil {
			return errors.New("用户名已存在")
		}
		user.Username = req.UserName
	}
	err = DB.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(id int) error {
	err := DB.Where("id=?", id).Delete(&User{}).Error
	if err != nil {
		return err
	}
	return nil
}
