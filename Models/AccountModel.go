package Models

import (
	"github.com/jinzhu/gorm"
)

type Account struct {
	gorm.Model
	Email        string   `json:"email"`
	Status       int    `json:"status"`
	Num 		int		 `json:"num"`
	Password	string   `json:"password"`
}

//默认表明为accounts    这里更改为account
func (Account) TableName() string {
	return "account"
}