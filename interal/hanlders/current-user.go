package hanlders

import (
	"github.com/gin-gonic/gin"
	"translate-system/serve"
)

func CurrentUserHandler(c *gin.Context) {
	username, err := serve.GetUsernameFromCookie(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "未登录"})
		return
	}
	c.JSON(200, gin.H{"username": username})
}
