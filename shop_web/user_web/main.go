package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"shop_web/user_web/global"
	"shop_web/user_web/initialize"
	"shop_web/user_web/utils"
	"shop_web/user_web/utils/register/consul"
	myValidator "shop_web/user_web/validator"
	"syscall"
)

func main() {
	//1.初始化logger
	initialize.InitLogger()

	//2.初始化配置文件
	initialize.InitConfig()

	//3.初始化routers
	Router := initialize.Routers()

	//4.初始化翻译
	if err  := initialize.InitTrans("zh"); err != nil{
		panic(err)
	}

	//5.初始化srv的连接
	initialize.InitSrvConn()

	viper.AutomaticEnv()
	debug := viper.GetBool("SHOP_DEBUG")
	if !debug{
		port, err := utils.GetFreePort()
		if err == nil{
			global.ServerConfig.Port = port
		}
	}

	//6.注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok{
		_ = v.RegisterValidation("mobile", myValidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "非法的手机号码！", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())

			return t
		})
	}

	//7.注册服务到consul并启动服务器
	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	if err := registerClient.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId); err != nil{
		zap.S().Panic("用户服务注册失败", err.Error())
	}
	zap.S().Debugf("启动服务器, 端口：%d", global.ServerConfig.Port)
	go func() {
		if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败", err.Error())
		}
	}()
	//6.等待接受终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := registerClient.DeRegister(serviceId); err != nil{
		zap.S().Panic("注销用户服务失败：", err.Error())
	}else {
		zap.S().Infof("注销用户服务成功")
	}
}
