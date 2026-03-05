package pusher

import (
	"bifrost/common/errorx"
	"bifrost/common/response"
	"bifrost/internal/logic/pusher"
	"bifrost/internal/model"
	"bifrost/svc"
	"github.com/gin-gonic/gin"
)

type BroadcastHandler struct {
	svcCtx *svc.ServerContext
	logic  *pusher.BroadcastLogic
}

func NewBroadcastHandler(svcCtx *svc.ServerContext) *BroadcastHandler {
	return &BroadcastHandler{
		svcCtx: svcCtx,
		logic:  pusher.NewBroadcastLogic(svcCtx),
	}
}

func (s *BroadcastHandler) Handle(c *gin.Context) {
	var req model.PushBroadcastReq

	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		response.Response(c, nil, errorx.NewGinBindParamError())
		return
	}

	err = s.logic.Logic(req)
	response.Response(c, nil, err)
}
