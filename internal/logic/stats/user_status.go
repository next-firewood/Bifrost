package stats

import (
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

func (l *UserStatusLogic) Logic(req model.UserStatusReq) (resp model.MetricsResp, err error) {
	client, err := l.svcCtx.Hub.ClientsMap(req.K, req.V)
	if err != nil {
		return resp, err
	}

	for _, c := range client {
		resp.List = append(resp.List, model.MetricsItem{
			Ud: c.Ud,
		})
	}

	resp.Total = int64(len(client))

	return resp, err
}
