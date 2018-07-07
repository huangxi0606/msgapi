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
	"MsgApi/Models"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
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
		max_addressee, _ := redis.Int(c.Do("HGET", "config:max_addressee", "value"))
		fmt.Print(max_addressee)
		MsgKeys, err := redis.Strings(c.Do("KEYS", "dispatch:msgTask:*"))
		if err != nil {
			log.Panic("get msg task key faild --- " + err.Error())
		}
		if len(MsgKeys) <1{
			context.JSON(http.StatusSeeOther,gin.H{
				"code":203,
				"message":"暂无任务",
			})
			return
		}
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
			if isActive == false {
				continue
			}
			hh, _ := redis.Strings(c.Do("SMEMBERS", activeMsgKey))
			if len(hh)<1{
				continue
			}
			if len(hh)<max_addressee{
				max_addressee =len(hh)
			}
			addressee := hh[:max_addressee]
			for _,val :=range addressee{
				c.Do("SREM",activeMsgKey,val)
			}
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
				"addressee": addressee,
			})
			if msgtaskid >0{
				break
			}
		}
	}


	func ReplyMsgTask(context *gin.Context){
		//t := time.Now().Format("2006-01-02 15:04:05")
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":'t',
		})
		return
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
		msg_task_id,ok :=context.GetQuery("msg_task_id")
		if !ok {
			//context.JSON(http.StatusSeeOther,gin.H{
			//	"code":203,
			//	"message":"msg_task_id is required",
			//})
			//return
		}
		fmt.Print(msg_task_id)
		status,ok :=context.GetQuery("status")
		if !ok {
			context.JSON(http.StatusSeeOther,gin.H{
				"code":203,
				"message":"status is required",
			})
			return
		}
		addressee,ok :=context.GetQuery("addressee")
		if !ok {

		}
		device_version,ok :=context.GetQuery("device_version")
		if !ok {

		}
		fmt.Print(device_version)
		log,ok :=context.GetQuery("log")
		if !ok {
			//context.JSON(http.StatusSeeOther,gin.H{
			//	"code":203,
			//	"message":"status is required",
			//})
			//return
		}
		fmt.Print(log)
		//链接redis
		c, err := redis.Dial("tcp", Config.REDIS_SERVER,redis.DialDatabase(Config.REDIS_DB))
		defer c.Close()
		if err != nil {
			panic("connect redis server faild -- " + err.Error())
		}
		//状态是字符串0
		Status, _ := strconv.Atoi(status)
		if Status ==0{
				task_key := "dispatch:msgTask:" + msg_task_id
				//msg, _ := redis.StringMap(c.Do("HGETALL", task_key))
				//current_num, _ := strconv.Atoi(msg["current_num"])
			c.Do("HINCRBY", task_key, "current_num", 1)
		}
		//if status == "0"{
		//	current_num =current_num+1;
		//
		//}
		//fmt.Print(status)
		deviceKey := "msg:device:" + machine
		c.Do("INCRBY", deviceKey, 1)
		db, err := gorm.Open("mysql", Config.MSQ)
		defer db.Close()
		if err != nil {
			panic("mysql db connect faild --- " + err.Error())
		}
		accounts:=Models.Account{}
		account :=db.Where(Models.Account{Email:email}).First(&accounts)
		account.Scan(&accounts)
		//fmt.Print(accounts.Status)
		//os.Exit(1)
		devices :=Models.Device{}
		db.Where(Models.Device{Sn:sn}).First(&devices).Scan(&devices)
		now :=time.Now().Format("2006-01-02 15:04:05") // 这是个奇葩,必须是这个时间点, 据说是go诞生之日, 记忆方法:6-1-2-3-4-5
		//# 2014-01-07 09:42:20
		fmt.Print(now)
		if Status >0{
			if Status ==1{
				devices.Status= 1
				//devices.failed_at
			}else{
				accounts.Status=1
			}
		}
		devices.Num += 1
		accounts.Num +=1
		db.Save(&devices)
		db.Save(&accounts)
		if len(addressee) >0{
			addressees :=Models.Addressee{}
			db.Where("number",addressee).First(&addressees).Scan(&addressees)
			addressees.Num +=1
			db.Save(&addressees)
		}
		//machiness :=Models.Machine{}
		//machiness.MsgTaskId = msgtaskid
		//machiness.Machine = machine
		//db.Save(&machiness)
		//mschiness =user1 := User{Name: "ScopeUser1", Age: 1}
		Statusy, _ := strconv.Atoi(status)
		if Statusy == 0{
			msglogs :=Models.MsgLog{}
			msglogs.Addressee =addressee
			msglogs.Status =status
			msglogs.Account_email =email
			msglogs.Deviced_sn =sn
			msglogs.Msg_task_id =msg_task_id
			msglogs.Log =log
			db.Save(&msglogs)
		}
		deviceLiveExpire, _ := redis.Int(c.Do("HGET", "config:app_device_active_expire", "value"))
		if deviceLiveExpire == 0 {
			deviceLiveExpire = 10
		}
		//更新设备活跃记录

		onlineDeviceLogKey := "app:online-device:" + machine + ":" + msg_task_id

		c.Do("HSET", onlineDeviceLogKey, "device_name", machine)
		c.Do("HSET", onlineDeviceLogKey, "task_id", msg_task_id)
		c.Do("EXPIRE", onlineDeviceLogKey, deviceLiveExpire*60)

		//更新设备状态
		deviceDataKey := "app:device:" + machine
		c.Do("HSET", deviceDataKey, "name", machine)
		c.Do("HSET", deviceDataKey, "version", device_version)
		c.Do("HSET", deviceDataKey, "status", 1) //标记设备为在线
		c.Do("HSET", deviceDataKey, "last_active_time", time.Now().Format("2006-01-02 15:04:05"))
		context.JSON(http.StatusOK,gin.H{
			"code":200,
			"message":"回执成功",
		})
	}
