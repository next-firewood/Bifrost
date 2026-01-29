package public

import (
	"bifrost/common/wss"
	"bifrost/internal/model"
	"bifrost/svc"
	"net/http"
)

type WebsocketConnLogic struct {
	svcCtx *svc.ServerContext
}

func NewWebsocketConnLogic(svcCtx *svc.ServerContext) *WebsocketConnLogic {
	return &WebsocketConnLogic{
		svcCtx: svcCtx,
	}
}

func (l WebsocketConnLogic) Logic(req model.WebsocketConnReq, w http.ResponseWriter, r *http.Request, connErr error) (err error) {
	return wss.ServeWs(req.UserData, l.svcCtx.Hub, w, r, connErr)
}
