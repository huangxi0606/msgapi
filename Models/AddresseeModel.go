package Models

import (
	"github.com/jinzhu/gorm"
)

type Addressee struct {
	gorm.Model
	Number        string   `json:"number"`
	Status       int    `json:"status"`
	Num 		int		 `json:"num"`
}

//默认表明为accounts    这里更改为account
func (Addressee) TableName() string {
	return "addressee"
}