package Handlers

//func GetPerson(c *gin.Context) {
//	id := c.Params.ByName(“id”)
//	var person Person
//	if err := db.Where(“id = ?”, id).First(&person).Error; err != nil {
//		c.AbortWithStatus(404)
//		fmt.Println(err)
//	} else {
//	c.JSON(200, person)
//	}
//}


import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"MsgApi/Config"

	"database/sql"
	"net/http"
)
type Huser struct {
	Id                int64
	Name        string
	HhProfile   Profile
	HhProfileId sql.NullInt64
}

type Profile struct {
	Id               int64
	Tel        string
}

func (Huser) TableName() string {
	return "huser"
}
func (Profile) TableName() string {
	return "profile"
}
func Get_Relation(ctx *gin.Context){
	db, err := gorm.Open("mysql", Config.MSQ)
	defer db.Close()
	if err != nil {
		log.Panic("mysql db connect faild --- " + err.Error())
	}
//BillingAddress:  Address{Address1: "Billing Address - Address 1"},
	user := Huser{
		Name:            "jinzhu",
		HhProfile: Profile{Tel: "Billing Address - Address 1"},
	}
	if err := db.Save(&user).Error; err != nil {
			ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data1": "cuowu",
		})
		return
	}
	var address1 Profile
	db.Model(&user).Related(&address1, "HhProfileId")
	db.Model(&user).Related(&address1, "BillingAddressId")
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data" :user,
	})








	//var huser Huser
	//var profile Profile
	////db.First(&huser).Scan(&huser)
	//db.Model(&huser).Related(&profile)
	////hhy :=db.First(&huser)
	//fmt.Print(profile.Tel)
	//os.Exit(1)
	//db.Model(&huser).Update("zz", 1)
	//db.First(&huser).Scan(&huser)
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"code": 200,
	//		"data1": huser,
	//	})
	//return
}
