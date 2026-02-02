package route

import (
	"bifrost/internal/handler/http/public"
	"bifrost/internal/handler/http/pusher"
	"bifrost/internal/handler/http/stats"
	"bifrost/svc"
	"github.com/gin-gonic/gin"
)

func RegisterHTTP(r *gin.Engine, svcCtx *svc.ServerContext) {
	r.Use(gin.Recovery())
	// WebSocket 升级接口（HTTP GET）
	r.GET("/ws-conn", public.NewWebsocketConnHandler(svcCtx).Handle)

	// API 路由组
	api := r.Group("/api")
	{
		// 推送接口
		push := api.Group("/push")
		{
			push.POST("/push", pusher.PushUser)                                             // 推送单个用户
			push.POST("/filter/broadcast", pusher.NewFilterBroadcastHandler(svcCtx).Handle) // 条件广播
			//push.POST("/broadcast", )        // 广播
		}

		// 状态查询接口
		status := api.Group("/status")
		{
			status.GET("/user", stats.NewUserStatusHandler(svcCtx).Handle) // 查询单个用户在线状态
			status.GET("/metrics", stats.NewMetricsHandler(svcCtx).Handle) // 全局连接指标
		}
	}
}
