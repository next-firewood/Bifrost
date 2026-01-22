package stats

import "github.com/gin-gonic/gin"

func UserStatus(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "UserStatus handler placeholder"})
}
