package Handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"MsgApi/Config"
	"github.com/garyburd/redigo/redis"
	_"github.com/garyburd/redigo/redisx"
	"math/rand"
	"strconv"
	"MsgApi/Models"
	"github.com/jinzhu/gorm"
)

//var count int
//err := db.Model(&Like{}).Where(&Like{Ip: ip, Ua: ua, Title: title}).Count(&count).Error
//if err != nil {
//	return false, err
//}
func GetMsgTask(context *gin.Context){
	machine,ok :=context.GetQuery("machine")
	if !ok {
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"machine is required",
		})
		return
	}
	fmt.Print(machine)
	email,ok :=context.GetQuery("email")
	if !ok {
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"email is required",
		})
		return
	}
	fmt.Print(email)
	sn,ok :=context.GetQuery("sn")
	if !ok {
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"sn is required",
		})
		return
	}
	fmt.Print(sn)
	//链接redis
	c, err := redis.Dial("tcp", Config.REDIS_SERVER,redis.DialDatabase(Config.REDIS_DB))
	defer c.Close()
	if err != nil {
		log.Panic("connect redis server faild --- " + err.Error())
	}
	db, err := gorm.Open("mysql", Config.MSQ)
	defer db.Close()
	if err != nil {
		log.Panic("mysql db connect faild --- " + err.Error())
	}
	accounts:=Models.Account{}
	account :=db.Where(Models.Account{Email:email}).First(&accounts)
	account.Scan(&accounts)
	//fmt.Print(accounts.Status)
	//os.Exit(1)
	if accounts.Status < 2{
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"该账号不能用",
		})
		return
	}
	devices :=Models.Device{}
	db.Where(Models.Device{Sn:sn}).First(&devices).Scan(&devices)
	if devices.Status <2{
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"该sn不能用",
		})
		return
	}
	var deviceTaskCountLimit int
	deviceTaskCountLimit, _ = redis.Int(c.Do("HGET", "config:msg_task_device_limit", "value"))
	deivceKey := "ad:device:" + machine
	finishedTaskCount, _ := redis.Int(c.Do("GET", deivceKey))
	if finishedTaskCount >= deviceTaskCountLimit {
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"该设备已完成今日任务数",
		})
		return
	}
	//MsgKeys, err := redis.Strings(c.Do("KEYS", "dispatch:msgTask:*"))
	//msgTasksKeys, err := redis.Strings(c.Do("KEYS", "dispatch:msgTask:*"))
	max_addressee, _ := redis.Int(c.Do("HGET", "config:config:max_addressee", "value"))
	fmt.Print(max_addressee)
	//找出可用的账号
	MsgKeys, err := redis.Strings(c.Do("KEYS", "dispatch:msgTask:*"))
	if err != nil {
		log.Panic("get ad task key faild --- " + err.Error())
	}
	//accountKeys, _ := redis.Strings(c.Do("KEYS", "dispatch:msgTask:*"))
	perms := rand.Perm(len(MsgKeys)) //打乱账号顺序
	var msgs []map[string]string
	fmt.Print(msgs)
	for _, perm := range perms {
		MsgKey := MsgKeys[perm]
		msg, _ := redis.StringMap(c.Do("HGETALL", MsgKey))
		msg_status, _ := strconv.Atoi(msg["status"])
		msg_enable, _ := strconv.Atoi(msg["enable"])
		current_num, _ := strconv.Atoi(msg["current_num"])
		target_num, _ := strconv.Atoi(msg["target_num"])
		max_device_num, _ := strconv.Atoi(msg["max_device_num"])
		msgtaskid,_ := strconv.Atoi(msg["id"])
		//检查账号状态
		if msg_status != 1 || msg_enable != 1 {
			continue
		}
		if current_num >= target_num {
			continue
		}
		machines :=Models.Machine{}
		var count int
		db.Where(Models.Machine{Machine:machine}).Find(&machines).Count(&count)
		if count >max_device_num{
			var num int
			db.Where(Models.Machine{Machine:machine,MsgTaskId:msgtaskid}).Find(&machines).Count(&num)
			if num ==0{
				continue
			}
		}
		if current_num +max_addressee >target_num{
			max_addressee =target_num-current_num
		}

		activeMsgKey := "msg:app-enable-account:" + msg["id"]
		isActive, _ := redis.Bool(c.Do("EXISTS", activeMsgKey))
		if isActive {
			continue
		}
		hh, _ := redis.Strings(c.Do("SMEMBERS", activeMsgKey))
		addressee := hh[:6]
		for _,val :=range addressee{
			c.Do("SREM",activeMsgKey,val)
		}
		//context.JSON(http.StatusSeeOther,gin.H{
		//	"code":200,
		//	"message":addressee,
		//	"id":msg["id"],
		//})
		//return
		machiness :=Models.Machine{}
		machiness.MsgTaskId = msgtaskid
		machiness.Machine = machine
		db.Save(&machiness)
		context.JSON(http.StatusOK,gin.H{
			"code":200,
			"msg_task_id" :msg["id"],
			"msg_name":msg["name"] ,
			"msg_content":msg["content"],
			"msg_urls":msg["urls"],
			"msg_link":msg["link"],
			"num":max_addressee,
			"addressee": devices.Wifu,
		})
		if msgtaskid >0{
			break
		}
	}

	}

	//func Test(ctx *gin.Context){
	//	c, err := redis.Dial("tcp", Config.REDIS_SERVER,redis.DialDatabase(Config.REDIS_DB))
	//	defer c.Close()
	//	if err != nil {
	//		log.Panic("get task api,connect redis server faild --- " + err.Error())
	//	}
	//	//c, err := redis.Dial("tcp", Config.REDIS_SERVER)
	//	//defer c.Close()
	//	//if err != nil {
	//	//	log.Panic("connect redis server faild --- " + err.Error())
	//	//}
	//	//任务
	//	MsgKeys, err := redis.Strings(c.Do("KEYS", "dispatch:msgTask:*"))
	//	if err != nil {
	//		log.Panic("get hour task key faild --- " + err.Error())
	//	}
	//	perms := rand.Perm(len(MsgKeys)) //打乱顺序
	//	var msgs []map[string]string
	//	fmt.Print(msgs)
	//	//ctx.JSON(http.StatusSeeOther,gin.H{
	//	//	"code":200,
	//	//	"message":len(MsgKeys),
	//	//})
	//	//return
	//
	//
	//	for _, perm := range perms {
	//		MsgKey := MsgKeys[perm]
	//		msg, _ := redis.StringMap(c.Do("HGETALL", MsgKey))
	//		activeMsgKey := "msg:app-enable-account:" + msg["id"]
	//		hh, _ := redis.Strings(c.Do("SMEMBERS", activeMsgKey))
	//		addressee := hh[:6]
	//		for _,val :=range addressee{
	//			c.Do("SREM",activeMsgKey,val)
	//		}
	//		ctx.JSON(http.StatusSeeOther,gin.H{
	//			"code":200,
	//			"message":addressee,
	//			"id":msg["id"],
	//		})
	//		return
	//
	//
	//		//fmt.Println(hh)
	//		//os.Exit(1)
	//	}
	//
	//}
