package Models

import "time"

type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	DeletedAt *time.Time   `json:"deleted_at"`
}
type BaseResponse struct {
	Status bool        `json:"status"`
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
}

type Response struct {
	Status bool        `json:"status"`
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}
type Count struct {
	number   int    `json:"number"`
}