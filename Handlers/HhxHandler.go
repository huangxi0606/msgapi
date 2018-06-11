package Handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"MsgApi/Config"
	"net/http"
)

type Hhx struct{
	FirstName        string   `json:"first_name"`
	LastName       string    `json:"last_name"`
}
//func (Hhx) TableName() string {
//	return "hhx"
//}
func Get_hhx(context *gin.Context){
	db, err := gorm.Open("mysql", Config.MSQ)
	defer db.Close()
	if err != nil {
		log.Panic("mysql db connect faild --- " + err.Error())
	}
	hhx := Hhx{FirstName: "Xi",LastName: "Huang"}
	db.Create(&hhx)
	context.JSON(http.StatusOK,gin.H{
		"code":200,
		"email": hhx,
	})

}
