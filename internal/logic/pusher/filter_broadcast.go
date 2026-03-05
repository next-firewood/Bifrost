package pusher

import (
	"bifrost/internal/model"
	"bifrost/svc"
)

type FilterBroadcastLogic struct {
	svcCtx *svc.ServerContext
	Fn     func(message string) error
}

func NewFilterBroadcastLogic(svcCtx *svc.ServerContext) *FilterBroadcastLogic {
	return &FilterBroadcastLogic{
		svcCtx: svcCtx,
	}
}

func (l *FilterBroadcastLogic) Logic(req model.PushBroadcastReq) (err error) {
	return l.Fn(req.Msg)
}
