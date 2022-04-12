package model

import (
	"errors"
	"micro-cloudStorage/file/service"
	"time"
)

type FileSystem struct {
	Id           int `json:"id"`
	UserId       int
	ParentDictId int
	FileName     string
	EncryptedKey string
	Ctime        time.Time
	Mtime        time.Time
	Atime        time.Time
	FileType     string
	FileSize     int
}

func CreateNewRoot() (int, error) {
	var file = &FileSystem{
		ParentDictId: 0,
		FileName:     "/",
		Ctime:        time.Now(),
		Mtime:        time.Now(),
		Atime:        time.Now(),
		FileType:     "d",
		FileSize:     0,
	}
	err := DB.Create(&file).Error
	if err != nil {
		return 0, err
	}
	return file.Id, nil
}

func GetAll(req *service.FileRequest) ([]*service.FileDetails, error) {
	var files []*FileSystem
	req.PageNum = req.GetPageSize() * (req.GetPageNum() - 1)
	err := DB.Where("parent_dict_id = ? and user_id = ?", req.ID, req.UserID).Offset(req.PageNum).Limit(req.PageSize).Find(&files).Error
	if err != nil {
		return nil, err
	}
	res := make([]*service.FileDetails, len(files))

	for i, file := range files {
		res[i] = new(service.FileDetails)
		res[i].ID = int32(file.Id)
		res[i].FileName = file.FileName
		res[i].FileType = file.FileType
		res[i].EncryptedKey = file.EncryptedKey
	}

	return res, nil
}

func GetParentId(req *service.FileRequest) (int, error) {
	var err error
	var file = FileSystem{}
	err = DB.Select("parent_dict_id").Where("id = ?", req.ID).Find(&file).Error
	if err != nil {
		return 0, errors.New("查询失败")
	}
	return file.ParentDictId, nil
}

func Count(req *service.FileRequest) (int, error) {
	var count int
	err := DB.Model(&FileSystem{}).Where("parent_dict_id = ?", req.ID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func AddFile(req *service.FileRequest) (int, error) {
	file := FileSystem{
		ParentDictId: int(req.ParentDictId),
		FileName:     req.FileName,
		UserId:       int(req.UserID),
		EncryptedKey: req.EncryptedKey,
		Ctime:        time.Now(),
		Mtime:        time.Now(),
		Atime:        time.Now(),
		FileType:     "-",
		FileSize:     int(req.Size),
	}

	err := DB.Save(&file).Error
	if err != nil {
		return 0, errors.New("保存失败")
	}

	return file.Id, nil
}

func AddDict(req *service.FileRequest) (int, error) {
	file := FileSystem{
		ParentDictId: int(req.ParentDictId),
		FileName:     req.FileName,
		UserId:       int(req.UserID),
		EncryptedKey: req.EncryptedKey,
		Ctime:        time.Now(),
		Mtime:        time.Now(),
		Atime:        time.Now(),
		FileType:     "d",
		FileSize:     0,
	}
	err := DB.Save(&file).Error
	return file.Id, err
}

func DeleteFile(id int) error {
	return DB.Where("id = ?", id).Delete(&FileSystem{}).Error
}

func GetFile(id int) (string, error) {
	var file = &FileSystem{}
	err := DB.Where("id = ?", id).Find(&file).Error
	if err != nil {
		return "", err
	}
	return file.FileName, nil
}
