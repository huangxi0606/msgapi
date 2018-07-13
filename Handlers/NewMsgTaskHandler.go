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
	func GetNewMsgTask(context *gin.Context){
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
		sn,ok :=context.GetQuery("serial")
		if !ok {
			context.JSON(http.StatusSeeOther,gin.H{
				"code":203,
				"message":"serial is required",
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
			max_adress,_ := strconv.Atoi(msg["max_adress"])

			if max_adress ==0{
				max_adress,_= redis.Int(c.Do("HGET", "config:max_addressee", "value"))
			}
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
			fmt.Print(max_adress)
			if current_num + max_adress >target_num{
				max_adress =target_num-current_num
			}

			var yy [] string
			//file_column,_ :=msg["file_column"]
			file_column,_ := msg["file_column"]


			if len(file_column) >0 {
				//有txt
				fmt.Print(file_column)
				len, _ := redis.Int(c.Do("llen", "msg:txt:" +msg["id"]))
				if len<1{
					continue
				}

				if len <max_adress{
					max_adress =len
				}
				address_key :="msg:txt:"+msg["id"]
				fmt.Print(address_key)
				var addressee = make([]string,0,max_adress)
				for a := 0; a < max_adress; a++ {
					hh, _ := redis.String(c.Do("lpop", address_key))
					addressee = append(addressee,hh)
				}
				yy =addressee
			} else{
				activeMsgKey := "msg:app-enable-account:" + msg["id"]
				//context.JSON(http.StatusSeeOther,gin.H{
				//	"code":789,
				//	"message":activeMsgKey,
				//})
				//return
				isActive, _ := redis.Bool(c.Do("EXISTS", activeMsgKey))
				if isActive == false {
					continue
				}
				hh, _ := redis.Strings(c.Do("SMEMBERS", activeMsgKey))
				if len(hh)<1{
					continue
				}
				if current_num + max_adress >target_num{
					max_adress =target_num-current_num
				}
				if len(hh)<max_adress{
					max_adress =len(hh)
				}
				addressee := hh[:max_adress]
				for _,val :=range addressee{
					c.Do("SREM",activeMsgKey,val)
				}
				yy =addressee
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
				"num":max_adress,
				"addressee": yy,
			})
			if msgtaskid >0{
				break
			}
		}
	}

	func ReplyNewMsgTask(context *gin.Context){
		//t := time.Now().Format("2006-01-02 15:04:05")
		//context.JSON(http.StatusSeeOther,gin.H{
		//	"code":203,
		//	"message":t,
		//})
		//return
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
		sn,ok :=context.GetQuery("serial")
		if !ok {
			context.JSON(http.StatusSeeOther,gin.H{
				"code":203,
				"message":"serial is required",
			})
			return
		}
		fmt.Print(sn)
		msg_task_id,ok :=context.GetQuery("msg_task_id")
		if !ok {
			context.JSON(http.StatusSeeOther,gin.H{
				"code":203,
				"message":"msg_task_id is required",
			})
			return
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
			context.JSON(http.StatusSeeOther,gin.H{
				"code":203,
				"message":"addressee is required",
			})
			return
		}

		msg_device_name,ok :=context.GetQuery("msg_device_name")
		if !ok {
			context.JSON(http.StatusSeeOther,gin.H{
				"code":203,
				"message":"msg_device_name is required",
			})
			return
		}
		device_version,ok :=context.GetQuery("device_version")
		if !ok {
			context.JSON(http.StatusSeeOther,gin.H{
				"code":203,
				"message":"device_version is required",
			})
			return
		}
		fmt.Print(device_version)
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
			c.Do("HINCRBY", task_key, "current_num", 1)
		}
		deviceKey := "msg:device:" + machine
		c.Do("INCRBY", deviceKey, 1)
		//设备和账号加1
		device := "machine:device:" + sn
		v, err := redis.Values(c.Do("HGETALL", device))
		if err != nil {
			panic(err)
		}
		if err := redis.ScanStruct(v, &p3); err != nil {
			panic(err)
		}
		c.Do("HSET", device, "num", p3.Num +1)

		account := "machine:device:" + email
		w, err := redis.Values(c.Do("HGETALL", account))
		if err != nil {
			panic(err)
		}
		if err := redis.ScanStruct(w, &p2); err != nil {
			panic(err)
		}
		c.Do("HSET", account, "num", p2.Num +1)
		db, err := gorm.Open("mysql", Config.MSQ)
		defer db.Close()
		if err != nil {
			panic("mysql db connect faild --- " + err.Error())
		}
		if len(addressee) >0{
			addressees :=Models.Addressee{}
			db.Where("number",addressee).First(&addressees).Scan(&addressees)
			addressees.Num +=1
			db.Save(&addressees)
		}

		log := "msg:log:" + msg_task_id+":"+addressee+":"+status
		c.Do("HSET", log, "status", status)
		c.Do("HSET", log, "msg_task_id", msg_task_id)
		c.Do("HSET", log, "deviced_sn", sn)
		c.Do("HSET", log, "account_email", email)
		c.Do("HSET", log, "msg_device_name", msg_device_name)
		c.Do("HSET", log, "account_email", email)
		c.Do("HSET", log, "msg_device_name", msg_device_name)
		c.Do("HSET", log, "addressee", addressee)
		c.Do("HSET", log, "current_at", time.Now().Format("2006-01-02 15:04:05"))
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
	type Log struct {
		Email     string `form:"email" json:"email" binding:"required"`
		Serial string `form:"serial" json:"serial" binding:"required"`
		Machine string `form:"machine" json:"machine" binding:"required"`
		Taskid string `form:"task_id" json:"task_id" binding:"required"`
		//Addt  []map[string]interface{} `form:"addressee" json:"addressee" binding:"required"`
		Addt  []map[string]string `form:"addressee" json:"addressee" binding:"required"`
	}
	func GetJson(context *gin.Context){
		c, err := redis.Dial("tcp", Config.REDIS_SERVER,redis.DialDatabase(Config.REDIS_DB))
		defer c.Close()
		if err != nil {
			panic("connect redis server faild -- " + err.Error())
		}
		// bind JSON数据
		var json Log
		if context.BindJSON(&json) == nil {
			msg_task_id :=json.Taskid
			hhy :=json.Addt
			machine :=json.Machine
			for _, value := range hhy{
				for k,v :=range value{
					//可能需要改为value['k']
					b,error := strconv.Atoi(v)
					if error != nil{
						fmt.Println("字符串转换成整数失败")
					}
					//k为收信人  b为状态
					task_key := "dispatch:msgTask:" + msg_task_id
					if b ==0{
						c.Do("HINCRBY", task_key, "current_num", 1)
					}
					if b ==2 {
						msg, _ := redis.StringMap(c.Do("HGETALL", task_key))
						file_column,_ := msg["file_column"]
						if len(file_column) >0{
							address_key := "msg:app-enable-account:" +msg_task_id
							_, err := redis.String(c.Do("SADD", address_key, k))
							if err != nil {
								log.Println("SADD failed:", err)
								return
							}
						}else{
							_, err = c.Do("rpush", "msg:txt:"+msg_task_id, k)
							if err != nil {
								fmt.Println("redis set failed:", err)
							}
						}
					}
					db, err := gorm.Open("mysql", Config.MSQ)
					defer db.Close()
					if err != nil {
						panic("mysql db connect faild --- " + err.Error())
					}
					if len(k) >0{
						addressees :=Models.Addressee{}
						db.Where("number",k).First(&addressees).Scan(&addressees)
						addressees.Num +=1
						db.Save(&addressees)
					}
					ll := "msg:log:" + msg_task_id+":"+k+":"+v
					c.Do("HSET", ll, "status", b)
					c.Do("HSET", ll, "msg_task_id", msg_task_id)
					c.Do("HSET", ll, "deviced_sn", json.Serial)
					c.Do("HSET", ll, "account_email", json.Email)
					c.Do("HSET", ll, "addressee", k)
					c.Do("HSET", ll, "current_at", time.Now().Format("2006-01-02 15:04:05"))
				}
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
			context.JSON(http.StatusOK,gin.H{
				"code":200,
				"message":"回执成功",
			})

		} else {
			context.JSON(http.StatusSeeOther,gin.H{
				"code":203,
				"message":"no",
			})
			return
			context.JSON(404, gin.H{"JSON=== status": "binding JSON error!"})
		}


	}
