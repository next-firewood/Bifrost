package public

import (
	"bifrost/common/response"
	"bifrost/internal/logic/public"
	"bifrost/internal/model"
	"bifrost/svc"
	"github.com/gin-gonic/gin"
)

type WebSocketConnHandler struct {
	svcCtx *svc.ServerContext
	logic  *public.WebsocketConnLogic
}

func NewWebsocketConnHandler(svcCtx *svc.ServerContext) *WebSocketConnHandler {
	return &WebSocketConnHandler{
		svcCtx: svcCtx,
		logic:  public.NewWebsocketConnLogic(svcCtx),
	}
}
func (s *WebSocketConnHandler) Handle(c *gin.Context) {
	token := c.GetHeader("Authorization")

	cla, err := s.svcCtx.Config.Auth.ValidateToken(token)
	if err != nil {
		return
	}

	err = s.logic.Logic(model.WebsocketConnReq{UserData: cla.Ud}, c.Writer, c.Request, err)
	response.Response(c, nil, err)
}
