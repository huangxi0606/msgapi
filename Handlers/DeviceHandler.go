package Handlers

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"log"
	"MsgApi/Config"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"github.com/garyburd/redigo/redis"
	"time"
)
var  p3 struct {
	id  string `redis:"id"`
	Serial string `redis:"sn"`
	Udid   string `redis:"udid"`
	BluetoothAddress string `redis:"bt"`
	Imei string `redis:"imei"`
	WiFiAddress   string `redis:"wifi"`
	Ecid string `redis:"ecid"`
	ProductType string `redis:"tp"`
	ModelNumber   string `redis:"nb"`
	RegionInfo string `redis:"reg"`
	Ethernet string `redis:"ethernet"`
	ICCID   string `redis:"iccid"`
	IMSI string `redis:"IMSI"`
	BasebandSerialNumber string `redis:"BasebandSerialNumber"`
	BasebandMasterKeyHash string `redis:"BasebandMasterKeyHash"`
	BasebandChipID   string `redis:"BasebandChipID"`
	Record_at string `redis:"Record_at"`
	Num int `redis:"num"`
}

func GetNewDevice(context *gin.Context){

	machine,ok :=context.GetQuery("machine")
	if !ok {
		//context.String(203, "machine is required")
		//return
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"machine is required",
		})
		return
	}
	fmt.Print(machine)

	status,ok :=context.GetQuery("status")
	if !ok {
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"status is required",
		})
		return
	}
	fmt.Print(status)
	tp,ok :=context.GetQuery("ProductType")
	if !ok {
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"ProductType is required",
		})
		return
	}
	fmt.Print(tp)
	//链接redis
	c, err := redis.Dial("tcp", Config.REDIS_SERVER,redis.DialDatabase(Config.REDIS_DB))
	defer c.Close()
	if err != nil {
		log.Panic("connect redis server faild --- " + err.Error())

	}
	values, err := redis.String(c.Do("lpop", "device:list:"+status+":"+tp))
	if err != nil{
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"合适的设备不存在",
		})
		return
	}
	fmt.Print(values)
	v, err := redis.Values(c.Do("HGETALL", "msg:device:"+values))
	if err != nil {
		panic(err)
	}
	if err := redis.ScanStruct(v, &p3); err != nil {
		panic(err)
	}

	c.Do("HSET", "old:device:"+p3.Serial, "now_time", time.Now().Format("2006-01-02 15:04:05"))
	c.Do("HSET", "old:device:"+p3.Serial, "sn", p3.Serial)
	c.Do("HSET", "old:device:"+p3.Serial, "status", status)
	context.JSON(http.StatusOK,gin.H{
		"code":200,
		"Serial" :p3.Serial,
		"Udid":p3.Udid ,
		"BluetoothAddress":p3.BluetoothAddress,
		"Imei" :p3.Imei,
		"WiFiAddress":p3.WiFiAddress ,
		"Ecid":p3.Ecid,
		"ProductType" :p3.ProductType,
		"ModelNumber":p3.ModelNumber ,
		"RegionInfo":p3.RegionInfo,
		"Ethernet" :p3.Ethernet,
		"ICCID":p3.ICCID ,
		"IMSI":p3.IMSI,
		"BasebandSerialNumber":p3.BasebandSerialNumber,
		"BasebandMasterKeyHash" :p3.BasebandMasterKeyHash,
		"BasebandChipID":p3.BasebandChipID ,
		"Record_at":p3.Record_at,
	})
	return
}

func ReplyNewDevice(context *gin.Context){
	//链接redis
	c, err := redis.Dial("tcp", Config.REDIS_SERVER,redis.DialDatabase(Config.REDIS_DB))
	defer c.Close()
	if err != nil {
		log.Panic("connect redis server faild --- " + err.Error())

	}
	status,ok :=context.GetQuery("status")
	if !ok {
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"status is required",
		})
		return
	}
	fmt.Print(status)

	tp,ok :=context.GetQuery("ProductType")
	if !ok {
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"ProductType is required",
		})
		return
	}
	fmt.Print(tp)
	serial,ok :=context.GetQuery("serial")
	if !ok {
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"serial is required",
		})
		return
	}
	fmt.Print(serial)
	v, err := redis.Values(c.Do("HGETALL", "msg:device:"+serial))
	if err != nil {
		panic(err)
	}
	if err := redis.ScanStruct(v, &p3); err != nil {
		panic(err)
	}
	c.Do("HSET", "msg:device:"+serial, "status", status)
	c.Do("HSET", "msg:device:"+serial, "updated_at", time.Now().Format("2006-01-02 15:04:05"))
	rec,ok :=context.GetQuery("rec")
	if ok {
		c.Do("HSET", "msg:device:"+serial, "Record_at", time.Now())
	}
	fmt.Print(rec)

	_, err = c.Do("rpush", "device:list:"+status+":"+tp, serial)
	if err != nil {
		fmt.Println("redis set failed:", err)
	}
	context.JSON(http.StatusOK,gin.H{
		"code":200,
		"message" :"回执成功",
	})
}


