package Handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"MsgApi/Config"
)

type Hht struct{
	Id               int64      `json:"id"`
	One        string   `json:"one"`
	Two       string    `json:"two"`
	Three       string    `json:"three"`
}
func (Hht) TableName() string {
	return "hht"
}
func Get_hhx(context *gin.Context){
	db, err := gorm.Open("mysql", Config.MSQ)
	defer db.Close()
	if err != nil {
		log.Panic("mysql db connect faild --- " + err.Error())
	}

	//var hht Hht
	//1.取出第一个数据（正确）
	//db.First(&hht).Scan(&hht)
	//2.保存数据  ok
	//var hht1 = Hht{One: "hebe", Two: "tien"}
	//var hht2 = Hht{One: "ella", Two: "chen"}
	//var hht3 = Hht{One: "selina", Two: "ren"}
	//var hht4 = Hht{One: "anpu", Two: "jiao"}
	//
	//db.Create(&hht1)
	//db.Create(&hht2)
	//db.Create(&hht3)
	//db.Create(&hht4)
	//3.list数据
	//var hht =[]Hht{}
	//db.Where("id >?", 0).Find(&hht).Scan(&hht)
	////fmt.Println(hhy)
	//for _, hh := range hht {
	//	context.JSON(http.StatusOK, gin.H{
	//		"code": 200,
	//		"One":  hh.One,
	//		"two":  hh.Two,
	//	})
	//}
	//return
	//4.更改某一列
	//db.Model(&hht).Update("three", "love")



}
