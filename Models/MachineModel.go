package Models

import "github.com/jinzhu/gorm"

type Machine struct {
	gorm.Model
	Machine        string   `json:"machine"`
	MsgTaskId       int    `json:"msg_task_id"`

}

//默认表明为accounts    这里更改为account
func (Machine) TableName() string {
	return "machine"
}
