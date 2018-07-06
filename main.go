package main

import (
	"github.com/gin-gonic/gin"
	."MsgApi/Handlers"
	"MsgApi/Middlewares"
)

func main(){
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(Middlewares.CORSMiddleware())//
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/v1/msg_task/get_account",GetAccount)
	r.GET("/v1/msg_task/get_device",GetDevice)
	r.GET("/v1/msg_task/get_msgtask",GetMsgTask)
	r.GET("/v1/msg_task/reply_msgtask",ReplyMsgTask)
	r.GET("/v1/msg_task/get_hhx",Get_hhx)
	r.GET("/v1/msg_task/get_relation",Get_Relation)
	r.GET("/v1/msg_task/new",GetNewAccount)
	r.GET("/v1/msg_task/test",Test)
	r.GET("/v1/msg_task/get_newaccount",GetNewAccount)
	r.POST("/v1/msg_task/reply_newaccount",ReplyNewAccount)
	r.GET("/v1/msg_task/get_newdevice",GetNewDevice)
	r.POST("/v1/msg_task/reply_newdevice",ReplyNewDevice)
	r.GET("/v1/msg_task/get_newmsgtask",GetNewMsgTask)
	r.Run(":9009")
}

