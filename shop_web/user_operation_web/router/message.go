package router

import (
	"github.com/gin-gonic/gin"
	"shop_web/user_operation_web/api/message"
	"shop_web/user_operation_web/middlewares"
)


func InitMessageRouter(Router *gin.RouterGroup) {
	MessageRouter := Router.Group("messages")
	{
		MessageRouter.GET("", middlewares.JWTAuth(), message.List) // 获取留言信息列表页
		MessageRouter.POST("", middlewares.JWTAuth(), message.New) //新建留言信息
		MessageRouter.DELETE("/:id", middlewares.JWTAuth(), message.Delete)//删除留言信息
	}
}
