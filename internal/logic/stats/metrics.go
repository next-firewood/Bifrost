package stats

import (
	"bifrost/internal/model"
	"bifrost/svc"
)

type MetricsLogic struct {
	svcCtx *svc.ServerContext
}

func NewMetricsLogic(svcCtx *svc.ServerContext) *MetricsLogic {
	return &MetricsLogic{
		svcCtx: svcCtx,
	}
}

func (l *MetricsLogic) Logic() (resp model.MetricsResp, err error) {
	for client, b := range l.svcCtx.Hub.Clients() {
		if b {
			resp.List = append(resp.List, model.MetricsItem{Ud: client.Ud})
			resp.Total++
		}
	}

	return resp, err
}
