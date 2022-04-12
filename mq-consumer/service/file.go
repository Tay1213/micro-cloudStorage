package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"micro-cloudStorage/mq-server/model"
	"micro-cloudStorage/mq-server/pkg/constant"
	"os"
	"time"
)

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

func CreateFile() {
	channel, err := model.MQ.Channel()
	if err != nil {
		panic(err)
	}
	//队列声明 durable：持久化， autoDelete：自动删除， exclusive：是否排他
	q, _ := channel.QueueDeclare("file_queue", true, false, false, false, nil)
	err = channel.Qos(1, 0, false)
	msg, err := channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	go func() {
		for m := range msg {
			f := File{}
			err = json.Unmarshal(m.Body, &f)
			if err != nil {
				panic(err)
			}

			file, err := os.OpenFile(constant.FileStoreRoot+f.FileName, os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				fmt.Println("打开文件失败")
				panic(err)
			}
			writer := bufio.NewWriter(file)
			nn, err := writer.Write(f.Data)
			if err != nil || nn != int(f.Size) {
				fmt.Println("写入文件失败")
				panic(err)
			}
			writer.Flush()
			file.Close()
			fileSystem := model.FileSystem{
				ParentDictId: f.ParentDictId,
				FileName:     f.FileName,
				UserId:       f.UserId,
				EncryptedKey: f.EncryptedKey,
				Ctime:        time.Now(),
				Mtime:        time.Now(),
				Atime:        time.Now(),
				FileType:     "-",
				FileSize:     f.Size,
			}
			_, err = model.AddFile(fileSystem)
			if err != nil {
				fmt.Println("添加文件失败")
				panic(err)
			}
			fmt.Println("写入成功")
			_ = m.Ack(false)
		}
	}()
}
