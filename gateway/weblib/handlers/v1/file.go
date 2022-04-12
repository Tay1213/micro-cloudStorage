package v1

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"micro-cloudStorage/gateway/pkg/app"
	"micro-cloudStorage/gateway/service"
	"mime/multipart"
	"strconv"
)

type FileForm struct {
	UploadFile   *multipart.FileHeader `form:"upload_file"`
	UserID       int                   `form:"user_id"`
	ParentDictId int                   `form:"parent_dict_id"`
	FileName     string                `form:"file_name"`
	EncryptedKey string                `form:"encrypted_key"`
}

func AddFile(c *gin.Context) {
	appG := app.Gin{C: c}

	var fileForm = &FileForm{}
	PanicIfFileError(c.Bind(&fileForm))
	file, err := fileForm.UploadFile.Open()
	if err != nil {
		fmt.Println("获取文件失败")
		appG.FailMsg(500, "上传文件失败")
	}
	data := make([]byte, fileForm.UploadFile.Size)
	file.Read(data)
	_, err = FileServices.AddFile(context.Background(), &service.FileRequest{
		ParentDictId: int32(fileForm.ParentDictId),
		UserID:       int32(fileForm.UserID),
		EncryptedKey: fileForm.EncryptedKey,
		FileName:     fileForm.UploadFile.Filename,
		Size:         int32(fileForm.UploadFile.Size),
		Data:         data,
	})
	PanicIfFileError(err)
	appG.Success(nil)
}

func AddDict(c *gin.Context) {
	appG := app.Gin{C: c}
	var req = service.FileRequest{}
	PanicIfFileError(c.Bind(&req))
	resp, err := FileServices.AddDict(context.Background(), &req)
	PanicIfFileError(err)
	appG.Success(resp)
}

func GetFiles(c *gin.Context) {
	appG := app.Gin{C: c}
	var fileReq service.FileRequest
	PanicIfFileError(c.Bind(&fileReq))
	fileResp, err := FileServices.GetFilesByParentDictId(context.Background(), &fileReq)
	PanicIfUserError(err)
	appG.Success(fileResp)
}

func GetFile(c *gin.Context) {
	appG := app.Gin{C: c}
	var fileReq service.FileRequest
	id, err := strconv.Atoi(c.Param("id"))
	PanicIfFileError(err)
	fileReq.ID = int32(id)
	fileResp, err := FileServices.GetFile(context.Background(), &fileReq)
	PanicIfUserError(err)
	appG.Success(fileResp)
}

func DeleteFile(c *gin.Context) {
	appG := app.Gin{C: c}
	var fileReq service.FileRequest
	PanicIfFileError(c.Bind(&fileReq))
	fileResp, err := FileServices.DeleteFile(context.Background(), &fileReq)
	PanicIfUserError(err)
	appG.Success(fileResp)
}
