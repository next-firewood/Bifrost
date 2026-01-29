package stats

import (
	"bifrost/common/errorx"
	"bifrost/internal/model"
	"bifrost/svc"
)

type UserStatusLogic struct {
	svcCtx *svc.ServerContext
}

func NewUserStatusLogic(svcCtx *svc.ServerContext) *UserStatusLogic {
	return &UserStatusLogic{
		svcCtx: svcCtx,
	}
}

func (l *UserStatusLogic) Logic(req model.UserStatusReq) (resp model.MetricsItem, err error) {
	client := l.svcCtx.Hub.ClientsMap(req.Ckv)

	if client == nil {
		return resp, errorx.BusinessErr("当前连接不存在")
	}

	resp.Ud = client.Ud

	return resp, err
}
