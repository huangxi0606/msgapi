package Handlers

import (
	"github.com/gin-gonic/gin"
	"database/sql"
	"github.com/jinzhu/gorm"
	"log"
	"MsgApi/Config"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)
type User struct {
	Id                int64
	Name              string `sql:"size:255"`
	BillingAddress    Address       // Embedded struct
	BillingAddressID  sql.NullInt64 // Embedded struct's foreign key
}
type Address struct {
	ID        int
	Address1  string
	Address2  string
}
func Test (ctx *gin.Context){
	db, err := gorm.Open("mysql", Config.Test)
	defer db.Close()
	if err != nil {
		log.Panic("mysql db connect faild --- " + err.Error())
	}
	//user := User{
	//	Name:            "www",
	//	BillingAddress:  Address{Address1: "7893"},
	//}
	//if err := db.Save(&user).Error; err != nil {
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"code": 200,
	//		"data1": "one",
	//	})
	//	return
	//}
	var user User
	user = User{}
	db.First(&user).Scan(&user)
	var address1 Address
	db.Model(&user).Related(&address1, "BillingAddressId")
	ctx.JSON(http.StatusOK, gin.H{
		"code": 202,
		"data" :address1.Address1,
	})
	//return
	//if address1.Address1 != "123456" {
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"code": 200,
	//		"data1": "three",
	//	})
	//	return
	//}
	//var user1 User
	//user1 = User{}
	//db.Model(&address1).Related(&user1, "BillingAddressID")
	//if db.NewRecord(user1) {
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"code": 200,
	//		"data1": "four",
	//	})
	//	return
	//}
	//需解决
	//yy := make([]User, 0)
	//rows, err := db.Table("users").Rows()
	//for rows.Next() {
	//	var user User
	//	rows.Scan(&user)
	//	yy = append(yy, user)
	//
	//}
	//ctx.JSON(http.StatusOK, gin.H{
	//	"code": 200,
	//	"data1": yy,
	//})
	//return
   //var user User
	//db.Table("user").Joins("join addresses on users.billing_address_id =addresses.id").Scan(&user)
	//ctx.JSON(http.StatusOK, gin.H{
	//			"code": 200,
	//			"data1": user.BillingAddress.Address1,
	//		})
	//		return
}

