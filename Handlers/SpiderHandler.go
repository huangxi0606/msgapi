package Handlers

import (
	"github.com/gin-gonic/gin"
)

func Spider(ctx *gin.Context){
//百度
//	resp, err := http.Get("http://www.baidu.com")
//	if err != nil {
//		fmt.Println("http get error.")
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		fmt.Println("http read error")
//		return
//	}
//	src := string(body)
//	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
//	src = re.ReplaceAllStringFunc(src, strings.ToLower)
//	//去除STYLE
//	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
//	src = re.ReplaceAllString(src, "")
//	//去除SCRIPT
//	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
//	src = re.ReplaceAllString(src, "")
//	//去除所有尖括号内的HTML代码，并换成换行符
//	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
//	src = re.ReplaceAllString(src, "\n")
//
//	//去除连续的换行符
//	re, _ = regexp.Compile("\\s{2,}")
//	src = re.ReplaceAllString(src, "\n")
//
//	ctx.JSON(http.StatusSeeOther,gin.H{
//		"code":203,
//		"src":strings.TrimSpace(src),
//	})
//	return
//糗事百科

}
