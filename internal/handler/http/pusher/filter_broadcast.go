package pusher

import (
	"bifrost/common/errorx"
	"bifrost/common/response"
	"bifrost/internal/logic/pusher"
	"bifrost/internal/model"
	"bifrost/svc"
	"fmt"
	"github.com/gin-gonic/gin"
)

type FilterBroadcastHandler struct {
	svcCtx *svc.ServerContext
	logic  *pusher.PushBroadcastLogic
}

func NewFilterBroadcastHandler(svcCtx *svc.ServerContext) *FilterBroadcastHandler {
	return &FilterBroadcastHandler{
		svcCtx: svcCtx,
		logic:  pusher.NewPushBroadcastLogic(svcCtx),
	}
}

func (s *FilterBroadcastHandler) Handle(c *gin.Context) {
	var req model.PushBroadcastReq

	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		response.Response(c, nil, errorx.NewGinBindParamError())
		return
	}

	s.logic.Fn, err = s.svcCtx.Hub.FilterSendFn(c)
	if err != nil {
		fmt.Println(err)
		response.Response(c, nil, err)
		return
	}

	err = s.logic.Logic(req)
	response.Response(c, nil, err)
}
