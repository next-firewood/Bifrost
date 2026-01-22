package pusher

import "github.com/gin-gonic/gin"

func PushBroadcast(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "PushBroadcast handler placeholder"})
}
