package sinking_sdk_go

import (
	"fmt"
	"github.com/spf13/viper"
	"testing"
	"time"
)

var server *Register
var serverConf *viper.Viper

func Test_main(t *testing.T) {
	//for i := 2080; i < 2080; i++ {
	//	server := New("106.52.89.187", "sinking-token", "test_token", "sinking-go-api-order", "sinking.go", "dev", "sinking-go-api", "127.0.0.1:"+strconv.Itoa(i))
	//	server.Listen()
	//}
	//time.Sleep(999999 * time.Second)
	server = New("106.52.89.187:80", "sinking-token", "test_token", "sinking.go", "dev")
	//注册并监听服务
	server.Register("sinking-go-api", "sinking-go-api-order", "106.52.89.187").UseService(map[string]string{
		"sinking-go-api": "sinking-go-api-order", //需要使用的服务
	}).Listen()
	//rpc调用服务
	//body, err := server.Rpc("sinking-go-api-order").Timeout(5).Method(http.MethodPost).ReTry(5).Call("/index/login", &Param{
	//	"user": "admin",
	//	"pwd":  "123456",
	//})
	//fmt.Println(body, err)
	//config拉取配置

	for {
		//fmt.Println(server.server)
		//fmt.Println(services)
		//body, err := server.Rpc("sinking-go-api-order").Timeout(5).Method(http.MethodPost).ReTry(5).Call("/index/login", &Param{
		//	"user": "admin",
		//	"pwd":  "123456",
		//})
		//fmt.Println(body, err)
		serverConf = server.Config("sinking-go-api").Name("database").Viper()
		fmt.Println(time.Now(), serverConf.Get("host"))
		time.Sleep(time.Second)
	}
	//server.Listen()
	//server2 := New("127.0.0.1:8888", "sinking-token", "test_token", "sinking-go-api-order", "sinking.go", "dev", "sinking-go-api", "127.0.0.1:8887")
	//server2.Listen()
	//server3 := New("127.0.0.1:8888", "sinking-token", "test_token", "sinking-go-api-order", "sinking.go", "dev", "sinking-go-api", "127.0.0.1:8886")
	//server3.Listen()
	//server4 := New("127.0.0.1:8888", "sinking-token", "test_token", "sinking-go-api-pay", "sinking.go", "dev", "sinking-go-api", "127.0.0.1:8885")
	//server4.Listen()
	//server5 := New("127.0.0.1:8888", "sinking-token", "test_token", "sinking-go-api-pay", "sinking.go", "dev", "sinking-go-api", "127.0.0.1:8884")
	//server5.Listen()
	//go func() {
	//	time.Sleep(5 * time.Second)
	//	for {
	//		data, err := server.
	//			Rpc("sinking-go-api-order").                          //调用服务名
	//			Method(http.MethodPost).                              //请求方式
	//			Header(map[string]string{"test": "test_data"}).       //请求头
	//			Timeout(5).                                           //超时熔断
	//			ReTry(5).                                             //最大重试次数
	//			Call("/api/service/register", &Param{"test": "test"}) //调用地址及内容
	//		fmt.Println(data, err)
	//		time.Sleep(time.Second)
	//	}
	//}()
	time.Sleep(999 * time.Second)
}
