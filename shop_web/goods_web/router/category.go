package router

import (
	"github.com/gin-gonic/gin"
	"shop_web/goods_web/api/category"
)

func InitCategoryRouter(Router *gin.RouterGroup) {
	CategoryRouter := Router.Group("categorys")
	{
		CategoryRouter.GET("", category.List)          //获取商品分类列表
		CategoryRouter.DELETE("/:id", category.Delete) //删除商品分类
		CategoryRouter.GET("/:id", category.Detail)    //获取分类详情
		CategoryRouter.POST("", category.New)          //新建商品分类
		CategoryRouter.PUT("/:id", category.Update)    //修改分类信息
	}
}
