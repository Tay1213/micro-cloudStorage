package router

import (
	"github.com/gin-gonic/gin"
	"micro-cloudStorage/fileServer/pkg/app"
	"path"
)

type File struct {
	UserId   int    `json:"user_id"`
	FilePath string `json:"file_path"`
}

func GetFile(context *gin.Context) {
	appG := app.Gin{C: context}
	id := context.GetInt("id")
	var file = File{}
	err := context.Bind(&file)
	if err != nil {
		appG.Fail()
		return
	}

	if id != file.UserId {
		appG.Unauthorized()
		return
	}

	var HttpContentType = map[string]string{
		".avi":  "video/avi",
		".mp3":  "audio/mp3",
		".mp4":  "video/mp4",
		".wmv":  "video/x-ms-wmv",
		".asf":  "video/x-ms-asf",
		".rm":   "application/vnd.rn-realmedia",
		".rmvb": "application/vnd.rn-realmedia-vbr",
		".mov":  "video/quicktime",
		".m4v":  "video/mp4",
		".flv":  "video/x-flv",
		".jpg":  "image/jpeg",
		".png":  "image/png",
		".pdf":  "application/pdf",
		".docx": "application/msword",
		".doc":  "application/msword",
	}
	filePath := file.FilePath
	//获取文件名称带后缀
	fileNameWithSuffix := path.Base(filePath)
	//获取文件的后缀
	fileType := path.Ext(fileNameWithSuffix)
	//获取文件类型对应的http ContentType 类型
	fileContentType := HttpContentType[fileType]
	if fileContentType == "" {
		fileContentType = "application/octet-stream"
	}
	context.Header("Content-Type", fileContentType)
	context.Header("Content-Disposition", "attachment")
	context.File(filePath)

}
