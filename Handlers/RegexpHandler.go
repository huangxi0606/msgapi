package Handlers

import (
	"github.com/gin-gonic/gin"
	"regexp"
	"net/http"
)

func Regexp(ctx *gin.Context)  {
	//判断是不是ip
	ip,ok :=ctx.GetQuery("ip")
	if !ok {
	}
	hx :=IsIP(ip)
	ctx.JSON(http.StatusSeeOther,gin.H{
		"code":203,
		"bool":hx,
	})
	return

}
func IsIP(ip string) (b bool) {
	if m, _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", ip); !m {
		return false
	}
	return true
}