package Handlers

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"log"
	"github.com/jinzhu/gorm"
	"MsgApi/Config"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"MsgApi/Models"
	"net/http"
)
func GetAccount(context *gin.Context){

	machine,ok :=context.GetQuery("machine")
	if !ok {
		//context.String(203, "machine is required")
		//return
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"machine is required",
		})
		return
	}
	fmt.Print(machine)
	db, err := gorm.Open("mysql", Config.MSQ)
	defer db.Close()
	if err != nil {
		log.Panic("mysql db connect faild --- " + err.Error())
	}
	var accounts =Models.Account{}
	db.Where(Models.Account{Status:0}).Order("num").First(&accounts).Scan(&accounts)
	db.Model(&accounts).Update("status", 2)
	//accounts.Status = 2
	if err := db.Save(&accounts).Error; err != nil {
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"Got errors when save accounts",
		})
		return
	}

	//db.Save(&accounts)
	//fmt.Println(accounts.Password)
	//os.Exit(1)
	context.JSON(http.StatusOK,gin.H{
		"code":200,
		"email": accounts.Email,
		"password" :accounts.Password,
	})
	//return
}

func GetDevice(context *gin.Context){
	machine,ok :=context.GetQuery("machine")
	if !ok {
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"machine is required",
		})
		return
	}
	fmt.Print(machine)
	ptype,ok :=context.GetQuery("ptype")
	if !ok {
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message": "type is required",
		})
		return
	}
	fmt.Print(ptype)
	db, err := gorm.Open("mysql", Config.MSQ)
	defer db.Close()
	if err != nil {
		log.Panic("mysql db connect faild --- " + err.Error())
	}
	//db.Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&result)
	var devices =Models.Device{}
	db.Where(Models.Device{Status:0,Type:ptype}).Order("num").First(&devices).Scan(&devices)
	if len(devices.Sn)<1{
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message": "没有合适的设备码可用",
		})
		return
	}
	//fmt.Print(devices)
	//os.Exit(1)
	db.Model(&devices).Update("status", 2)
	//devices.Status = 2
	if err := db.Save(&devices).Error; err != nil {
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"Got errors when save devices",
		})
		return
	}

	//db.Save(devices)
	//db.Save(&devices)
	context.JSON(http.StatusOK,gin.H{
		"code":200,
		"Udid" :devices.Udid,
		"Status":devices.Status ,
		"Num":devices.Num,
		"Sn":devices.Sn,
		"Imei":devices.Imei,
		"Bt":devices.Bt,
		"Wifu": devices.Wifu,
		"Ecid":  devices.Ecid,
		"Tp" :   devices.Tp,
		"Nb" :   devices.Nb,
		"Reg" :  devices.Reg,
		"Ethernet": devices.Ethernet,
		"ICCID":  devices.ICCID,
		"Type": devices.Type,
	})
	//return

}
