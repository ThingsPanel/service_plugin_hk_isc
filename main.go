package main

import (
	"fmt"
	"log"
	httpclient "service_hk_isc/http_client"
	httpservice "service_hk_isc/http_service"
	"service_hk_isc/mqtt"
	"service_hk_isc/services"
	"strings"

	"github.com/spf13/viper"
)

func main() {
	conf()
	LogInIt()
	log.Println("Starting the application...")
	// 启动mqtt客户端
	mqtt.InitClient()
	// 启动http客户端
	httpclient.Init()
	// 启动服务
	startService()
	// 启动http服务
	httpservice.Init()
	select {}
}

func startService() {
	server := NewServer()

	// 添加 HkIsc 服务
	server.AddService(services.NewHkIsc())
	// 添加其他服务...
	// server.AddService(newService1())
	// server.AddService(newService2())

	// 运行所有服务
	if err := server.Run(); err != nil {
		fmt.Printf("服务器运行出错: %v\n", err)
	}

	// 在程序结束时关闭所有服务
	defer server.Close()
}
func conf() {
	log.Println("加载配置文件...")
	// 设置环境变量前缀
	viper.SetEnvPrefix("service_hk_isc")
	// 使 Viper 能够读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")
	viper.SetConfigFile("./config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("加载配置文件完成...")
}
