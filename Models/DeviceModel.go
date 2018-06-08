package Models

import (
	"github.com/jinzhu/gorm"
)

//type Model struct {
//	ID        uint `gorm:"primary_key"`
//	CreatedAt time.Time  `json:"created_at"`
//	UpdatedAt time.Time   `json:"updated_at"`
//	DeletedAt *time.Time   `json:"deleted_at"`
//}

type Device struct {
	gorm.Model
	Udid        string   `json:"udid"`
	Status       int    `json:"status"`
	Num 		int		 `json:"num"`
	Sn			string   `json:"sn"`
	Imei        string   `json:"imei"`
	Bt			string   `json:"bt"`
	Wifu        string   `json:"wifu"`
	Ecid        string   `json:"ecid"`
	Tp        	string   `json:"tp"`
	Nb        	string   `json:"nb"`
	Reg         string   `json:"reg"`
	Ethernet    string   `json:"Ethernet"`
	ICCID       string   `json:"ICCID"`
	Type        string   `json:"type"`
}

//默认表明为accounts    这里更改为account
func (Device) TableName() string {
	return "device"
}