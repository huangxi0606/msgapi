package Handlers

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"fmt"
	"os"
	"log"
	"io"
)
func UploadPic(context *gin.Context){
	// 注意此处的文件名和client处的应该是一样的
	file, header , err := context.Request.FormFile("uploadFile")
	filename := header.Filename
	fmt.Println(header.Filename)
	// 创建临时接收文件(ke 更换文件名字)生成的文件在更目录下
	out, err := os.Create("copy_"+filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	// Copy数据
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	context.JSON(http.StatusOK,gin.H{
		"code":200,
		"message":"upload file success",

	})
	//return
}


