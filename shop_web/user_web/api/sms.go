package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"net/http"
	"shop_web/user_web/forms"
	"shop_web/user_web/global"
	"strings"
	"time"
)

func GenerateSmsCode(length int) string {
	//生成length长度的验证码
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < length; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func SendSms(ctx *gin.Context) {
	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil{
		HandleValidatorError(ctx, err)
		return
	}

	//client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", global.ServerConfig.AliSmsInfo.ApiKey, global.ServerConfig.AliSmsInfo.ApiSecret)
	//if err != nil {
	//	panic(err)
	//}
	smsCode := GenerateSmsCode(6)
	//request := requests.NewCommonRequest()
	//request.Method = "POST"
	//request.Scheme = "https"
	//request.Domain = "dysmsapi.aliyuncs.com"
	//request.Version = "2017-05-25"
	//request.ApiName = "SendSms"
	//request.QueryParams["RegionId"] = "cn-beijing"
	//request.QueryParams["PhoneNumbers"] = sendSmsForm.Mobile//手机号
	//request.QueryParams["SignName"] = "生鲜商城"//阿里云验证过的项目名
	//request.QueryParams["TemplateCode"] = "SMS_154950909"//阿里云的短信模板号
	//request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}"//短信模板中的验证码内容
	//response, err := client.ProcessCommonRequest(request)
	//fmt.Print(client.DoAction(request, response))
	//if err != nil {
	//	fmt.Print(err.Error())
	//}
	//将验证码与手机号码绑定并保存起来
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	rdb.Set(context.Background(), sendSmsForm.Mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second)

	ctx.JSON(http.StatusOK, gin.H{
		"msg":"短信发送成功",
	})
}
