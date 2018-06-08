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
	db.Where(Models.Account{Status:0}).First(&accounts).Scan(&accounts)
	accounts.Status = 2
	db.Save(&accounts)
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
		context.JSON(203,gin.H{
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
	var devices =Models.Device{}
	device :=db.Where(Models.Device{Status:0,Type:ptype}).First(&devices)
	device.Scan(&devices)
	fmt.Println(device)
	devices.Status = 2
	db.Save(&devices)
	//os.Exit(1)
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
