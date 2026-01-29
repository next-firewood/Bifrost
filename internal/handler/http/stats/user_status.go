package stats

import (
	"bifrost/common/errorx"
	"bifrost/common/response"
	"bifrost/internal/logic/stats"
	"bifrost/internal/model"
	"bifrost/svc"
	"github.com/gin-gonic/gin"
)

type UserStatusHandler struct {
	logic  *stats.UserStatusLogic
	svcCtx *svc.ServerContext
}

func NewUserStatusHandler(svcCtx *svc.ServerContext) *UserStatusHandler {
	return &UserStatusHandler{
		logic:  stats.NewUserStatusLogic(svcCtx),
		svcCtx: svcCtx,
	}
}

func (s *UserStatusHandler) Handle(c *gin.Context) {
	var req model.UserStatusReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Response(c, nil, errorx.NewGinBindParamError())
		return
	}

	resp, err := s.logic.Logic(req)
	response.Response(c, resp, err)
}
