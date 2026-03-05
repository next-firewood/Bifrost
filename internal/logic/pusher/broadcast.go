package pusher

import (
	"bifrost/common/errorx"
	"bifrost/internal/model"
	"bifrost/svc"
)

type BroadcastLogic struct {
	svcCtx *svc.ServerContext
}

func NewBroadcastLogic(svcCtx *svc.ServerContext) *BroadcastLogic {
	return &BroadcastLogic{
		svcCtx: svcCtx,
	}
}

func (l *BroadcastLogic) Logic(req model.PushBroadcastReq) (err error) {
	if len(req.Msg) <= 0 {
		return errorx.BusinessErr("消息不得为空")
	}
	for client, b := range l.svcCtx.Hub.Clients() {
		if b {
			client.Send <- []byte(req.Msg)
		}
	}

	return
}
