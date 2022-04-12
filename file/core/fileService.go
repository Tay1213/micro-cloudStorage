package core

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"micro-cloudStorage/file/model"
	"micro-cloudStorage/file/pkg/constant"
	"micro-cloudStorage/file/service"
	"os"
	"time"
)

type FileService struct {
}

type File struct {
	Id           int `json:"id"`
	UserId       int
	ParentDictId int
	FileName     string
	EncryptedKey string
	Ctime        time.Time
	Mtime        time.Time
	Atime        time.Time
	FileType     string
	Data         []byte
	Size         int
}

func (*FileService) CreateNewFileRoot(ctx context.Context, req *service.FileRequest, resp *service.FileResponse) error {
	rootId, err := model.CreateNewRoot()
	if err != nil {
		return err
	}
	resp.RootDictID = int32(rootId)
	return nil
}
func (*FileService) GetFilesByParentDictId(ctx context.Context, req *service.FileRequest, resp *service.FileResponse) error {
	count, err := model.Count(req)
	if err != nil {
		return err
	}
	files, err := model.GetAll(req)
	if err != nil {
		return err
	}
	parentId, err := model.GetParentId(req)
	if err != nil {
		return err
	}
	resp.TotalFile = int32(count)
	resp.ParentId = int32(parentId)
	resp.Files = files
	return nil
}
func (*FileService) AddFile(ctx context.Context, req *service.FileRequest, resp *service.FileResponse) error {

	f := File{
		Id:           int(req.ID),
		UserId:       int(req.UserID),
		ParentDictId: int(req.ParentDictId),
		FileName:     req.FileName,
		EncryptedKey: req.EncryptedKey,
		Ctime:        time.Now(),
		Mtime:        time.Now(),
		Atime:        time.Now(),
		FileType:     "-",
		Data:         req.Data,
		Size:         int(req.Size),
	}

	fileJson, err := json.Marshal(&f)
	if err != nil {
		return err
	}

	channel, err := model.MQ.Channel()
	if err != nil {
		return err
	}

	q, err := channel.QueueDeclare("file_create", true, false, false, false, nil)
	err = channel.Publish("", q.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         fileJson,
	})
	if err != nil {
		return err
	}
	file, err := os.OpenFile(constant.FileStoreRoot+req.FileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("打开文件失败")
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	nn, err := writer.Write(req.Data)
	if err != nil || nn != int(req.Size) {
		fmt.Println("写入文件失败")
		return err
	}
	writer.Flush()

	_, err = model.AddFile(req)
	if err != nil {
		fmt.Println("添加文件失败")
		return err
	}
	return nil
}

func (*FileService) AddDict(ctx context.Context, req *service.FileRequest, resp *service.FileResponse) error {
	id, err := model.AddDict(req)
	if err != nil {
		return err
	}
	resp.ID = int32(id)
	return nil
}

func (*FileService) DeleteFile(ctx context.Context, req *service.FileRequest, resp *service.FileResponse) error {
	return model.DeleteFile(int(req.ID))
}
func (*FileService) GetFile(ctx context.Context, req *service.FileRequest, resp *service.FileResponse) error {
	fileName, err := model.GetFile(int(req.ID))
	if err != nil {
		return err
	}
	resp.FileAddr = constant.FileStoreRoot + fileName
	return nil
}
