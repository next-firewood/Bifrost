package pusher

import "github.com/gin-gonic/gin"

func PushUser(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "PushUser handler placeholder"})
}
