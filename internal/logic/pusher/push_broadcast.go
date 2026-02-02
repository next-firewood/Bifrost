package pusher

import (
	"bifrost/internal/model"
	"bifrost/svc"
)

type PushBroadcastLogic struct {
	svcCtx *svc.ServerContext
	Fn     func(message string) error
}

func NewPushBroadcastLogic(svcCtx *svc.ServerContext) *PushBroadcastLogic {
	return &PushBroadcastLogic{
		svcCtx: svcCtx,
	}
}

func (l *PushBroadcastLogic) Logic(req model.PushBroadcastReq) (err error) {
	return l.Fn(req.Msg)
}
