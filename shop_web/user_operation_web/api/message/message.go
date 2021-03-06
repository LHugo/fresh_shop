package message

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"shop_web/user_operation_web/api"
	"shop_web/user_operation_web/forms"
	"shop_web/user_operation_web/global"
	"shop_web/user_operation_web/models"
	"shop_web/user_operation_web/proto"
	"strconv"
)

func List(ctx *gin.Context) {
	request := &proto.MessageRequest{}

	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}

	rsp, err := global.MessageClient.MessageList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("获取留言失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := map[string]interface{}{
		"total": rsp.Total,
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["user_id"] = value.UserId
		reMap["type"] = value.MessageType
		reMap["subject"] = value.Subject
		reMap["text"] = value.Text
		reMap["file"] = value.File

		result = append(result, reMap)
	}
	reMap["data"] = result

	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")

	messageForm := forms.MessageForm{}
	if err := ctx.ShouldBindJSON(&messageForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.MessageClient.CreateMessage(context.Background(), &proto.MessageRequest{
		UserId: int32(userId.(uint)),
		MessageType: messageForm.MessageType,
		Subject: messageForm.Subject,
		Text: messageForm.Text,
		File: messageForm.File,
	})

	if err != nil {
		zap.S().Errorw("添加留言失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":rsp.Id,
	})
}

func Delete(ctx *gin.Context){
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	userId, _ := ctx.Get("userId")
	_, err = global.MessageClient.DeleteMessage(context.Background(), &proto.MessageRequest{
		UserId:  int32(userId.(uint)),
		Id: int32(i),
	})
	if err != nil {
		zap.S().Errorw("删除留言信息失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":"删除留言成功",
	})
}