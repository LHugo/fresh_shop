package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop_web/user_web/global"
	"shop_web/user_web/proto"
)


func InitSrvConn(){
	//集成consul服务中心注册功能和负载均衡功能的grpc拨号连接srv
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		)
	if err != nil{
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}

//func InitSrvConn() {
//	//从注册中心获取到用户服务的信息
//	cfg := api.DefaultConfig()
//	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
//	userSrvHost := ""
//	userSrvPort := 0
//	client, err := api.NewClient(cfg)
//	if err != nil {
//		panic(err)
//	}
//	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerConfig.UserSrvInfo.Name))
//	if err != nil{
//		panic(err)
//	}
//	for _, value := range data{
//		userSrvHost = value.Address
//		userSrvPort = value.Port
//		break
//	}
//	if userSrvHost == ""{
//		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
//		return
//	}
//
//	//拨号连接用户grpc服务器
//	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
//	if err != nil {
//		zap.S().Errorw("连接【用户服务】失败",
//			"msg", err.Error(),
//		)
//	}
//	//生成grpc的client并调用接口
//	userSrvClient := proto.NewUserClient(userConn)
//	global.UserSrvClient = userSrvClient
//
//}
