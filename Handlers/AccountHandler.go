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
var  p2 struct {
	id  string `redis:"id"`
	Email string `redis:"email"`
	Password   string `redis:"password"`
	Cert string `redis:"cert"`
	Num int `redis:"num"`
}

//重写获取账号
func GetNewAccount(context *gin.Context){

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
	//链接redis
	c, err := redis.Dial("tcp", Config.REDIS_SERVER,redis.DialDatabase(Config.REDIS_DB))
	defer c.Close()
	if err != nil {
		log.Panic("connect redis server faild --- " + err.Error())

	}
	values, err := redis.String(c.Do("lpop", "account:list:"+status))
	if err != nil{
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"合适的账号不存在",
		})
		return
	}
	fmt.Print(values)
	v, err := redis.Values(c.Do("HGETALL", "msg:account:"+values))
	if err != nil {
		panic(err)
	}
	if err := redis.ScanStruct(v, &p2); err != nil {
		panic(err)
	}

	c.Do("HSET", "old:account:"+p2.Email, "now_time", time.Now().Format("2006-01-02 15:04:05"))
	c.Do("HSET", "old:account:"+p2.Email, "email", p2.Email)
	c.Do("HSET", "old:account:"+p2.Email, "status", status)
	context.JSON(http.StatusOK,gin.H{
		"code":200,
		"email" :p2.Email,
		"password":p2.Password ,
		"cert":p2.Cert,
	})
	return
}

func ReplyNewAccount(context *gin.Context){
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

	cert,ok :=context.GetQuery("cert")
	if !ok {
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"cert is required",
		})
		return
	}
	fmt.Print(cert)
	email,ok :=context.GetQuery("email")
	if !ok {
		context.JSON(http.StatusSeeOther,gin.H{
			"code":203,
			"message":"email is required",
		})
		return
	}
	fmt.Print(email)

	v, err := redis.Values(c.Do("HGETALL", "msg:account:"+email))
	if err != nil {
		panic(err)
	}
	if err := redis.ScanStruct(v, &p2); err != nil {
		panic(err)
	}

	c.Do("HSET", "msg:account:"+email, "cert", cert)
	c.Do("HSET", "msg:account:"+email, "status", status)
	c.Do("HSET", "msg:account:"+email, "updated_at", time.Now().Format("2006-01-02 15:04:05"))
	rec,ok :=context.GetQuery("rec")
	if ok {
		c.Do("HSET", "msg:account:"+email, "Record_at", time.Now())
	}
	fmt.Print(rec)

	_, err = c.Do("rpush", "account:list:"+status, email)
	if err != nil {
		fmt.Println("redis set failed:", err)
	}
	context.JSON(http.StatusOK,gin.H{
		"code":200,
		"message" :"回执成功",

	})

}


