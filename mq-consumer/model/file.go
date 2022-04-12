package model

import (
	"errors"
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

func AddFile(f FileSystem) (int, error) {

	err := DB.Save(&f).Error
	if err != nil {
		return 0, errors.New("保存失败")
	}

	return f.Id, nil
}
