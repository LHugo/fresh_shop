package router

import (
	"github.com/gin-gonic/gin"
	"shop_web/user_operation_web/api/favourite"
	"shop_web/user_operation_web/middlewares"
)

func InitFavRouter(Router *gin.RouterGroup) {
	UserFavRouter := Router.Group("favs")
	{
		UserFavRouter.DELETE("/:id", middlewares.JWTAuth(), favourite.Delete) // 删除收藏记录
		UserFavRouter.GET("/:id", middlewares.JWTAuth(), favourite.Detail)    // 查询收藏记录
		UserFavRouter.POST("", middlewares.JWTAuth(), favourite.New)          //新建收藏记录
		UserFavRouter.GET("", middlewares.JWTAuth(), favourite.List)          //获取当前用户的收藏
	}
}
