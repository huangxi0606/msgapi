package Handlers

import (
	"github.com/gin-gonic/gin"
)

//type Server struct {
//	ServerName string
//	ServerIP   string
//}

//type Serverslice struct {
//	Servers []Server
//}

//解析json数据
func Analysis(context *gin.Context){
	//yy,ok:=context.GetQuery("yy")
	//if !ok {
	//}
	//fmt.Print(yy)
//1.
//	var s Serverslice
//	str := `{"servers":[{"serverName":"Shanghai_VPN","serverIP":"127.0.0.1"},{"serverName":"Beijing_VPN","serverIP":"127.0.0.2"}]}`
//	json.Unmarshal([]byte(str), &s)
//	fmt.Println(s)
	//{[{Shanghai_VPN 127.0.0.1} {Beijing_VPN 127.0.0.2}]}

	//2.在json数据格式未知的情况下
	//b := []byte(`{"Name":"Wednesday","Age":16,"Parents":["hhy","hhx"]}`)
	//var f interface{}
	//err := json.Unmarshal(b, &f)
	//fmt.Print(err)
	//m := f.(map[string]interface{})
	//for k, v := range m {
	//	switch vv := v.(type) {
	//	case string:
	//		fmt.Println(k, "is string", vv)
	//	case int:
	//		fmt.Println(k, "is int", vv)
	//	case float64:
	//		fmt.Println(k,"is float64",vv)
	//	case []interface{}:
	//		fmt.Println(k, "is an array:")
	//		for i, u := range vv {
	//			fmt.Println(i, u)
	//		}
	//	default:
	//		fmt.Println(k, "is of a type I don't know how to handle")
	//	}
	//}

}
