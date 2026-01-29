package stats

import (
	"bifrost/common/response"
	"bifrost/internal/logic/stats"
	"bifrost/svc"
	"github.com/gin-gonic/gin"
)

type MetricsHandler struct {
	svcCtx *svc.ServerContext
	logic  *stats.MetricsLogic
}

func NewMetricsHandler(svcCtx *svc.ServerContext) *MetricsHandler {
	return &MetricsHandler{
		svcCtx: svcCtx,
		logic:  stats.NewMetricsLogic(svcCtx),
	}
}

func (s *MetricsHandler) Handle(c *gin.Context) {
	resp, err := s.logic.Logic()
	response.Response(c, resp, err)
}
