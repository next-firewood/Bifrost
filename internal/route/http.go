package route

import (
	"bifrost/internal/handler/http/public"
	"bifrost/internal/handler/http/pusher"
	"bifrost/internal/handler/http/stats"
	"github.com/gin-gonic/gin"
)

func RegisterHTTP(r *gin.Engine) {
	r.Use(gin.Recovery())
	// WebSocket 升级接口（HTTP GET）
	r.GET("/ws-conn", public.WebSocketConn)

	// API 路由组
	api := r.Group("/api")
	{
		// 推送接口
		push := api.Group("/push")
		{
			push.POST("/push", pusher.PushUser)           // 推送单个用户
			push.POST("/broadcast", pusher.PushBroadcast) // 广播
		}

		// 状态查询接口
		status := api.Group("/status")
		{
			status.GET("/user/:user_id", stats.UserStatus) // 查询单个用户在线状态
			status.GET("/metrics", stats.Metrics)          // 全局连接指标
		}
	}
}
