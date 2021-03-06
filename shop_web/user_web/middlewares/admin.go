package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_web/user_web/models"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		if currentUser.AuthorityId == 1{
			c.JSON(http.StatusForbidden, gin.H{
				"msg":"无权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}