package Handlers

import (
	"fmt"
	"github.com/robfig/cron"
	"log"
	"github.com/gin-gonic/gin"
)

//1）星号(*)
//表示 cron 表达式能匹配该字段的所有值。如在第5个字段使用星号(month)，表示每个月
//
//2）斜线(/)
//表示增长间隔，如第1个字段(minutes) 值是 3-59/15，表示每小时的第3分钟开始执行一次，之后每隔 15 分钟执行一次（即 3、18、33、48 这些时间点执行），这里也可以表示为：3/15
//
//3）逗号(,)
//用于枚举值，如第6个字段值是 MON,WED,FRI，表示 星期一、三、五 执行
//
//4）连字号(-)
//表示一个范围，如第3个字段的值为 9-17 表示 9am 到 5pm 直接每个小时（包括9和17）
//
//5）问号(?)
//只用于日(Day of month)和星期(Day of week)，\表示不指定值，可以用于代替 *
//1.1个定时任务
//func Cron(ctx *gin.Context)  {
	//i := 0
	//c := cron.New()
	//spec := "*/5 * * * * ?"
	//c.AddFunc(spec, func() {
	//	i++
	//	log.Println("cron running:", i)
	//})
	//c.Start()
	//
	//select{}

//}
//2.多个定时任务
type TestJob struct {
}

func (this TestJob)Run() {
	fmt.Println("testJob1...")
}

type Test2Job struct {
}

func (this Test2Job)Run() {
	fmt.Println("testJob2...")
}

//启动多个任务
func Cron(ctx *gin.Context) {
	i := 0
	c := cron.New()
	//AddFunc
	spec := "*/5 * * * * ?"
	c.AddFunc(spec, func() {
		i++
		log.Println("cron running:", i)
	})
	//AddJob方法
	c.AddJob(spec, TestJob{})
	c.AddJob(spec, Test2Job{})
	//启动计划任务
	c.Start()
	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	defer c.Stop()
	select{}
}
