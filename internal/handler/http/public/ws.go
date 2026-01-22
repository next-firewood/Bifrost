package public

import "github.com/gin-gonic/gin"

func WebSocketConn(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "ws handler placeholder"})
}
