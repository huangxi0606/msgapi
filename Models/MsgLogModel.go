package Models

import (
	"github.com/jinzhu/gorm"
)


type MsgLog struct {
	gorm.Model
	Msg_task_id        string   `json:"msg_task_id"`
	Deviced_sn			string   `json:"deviced_sn"`
	Account_email        string   `json:"account_email"`
	Msg_device_name			string   `json:"msg_device_name"`
	Status        string   `json:"status"`
	Log        string   `json:"log"`
	Addressee        	string   `json:"addressee"`
}

//默认表明为accounts    这里更改为account
func (MsgLog) TableName() string {
	return "msg_task_log"
}