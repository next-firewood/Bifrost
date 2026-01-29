package model

import "bifrost/common/jwtx"

type WebsocketConnReq struct {
	UserData map[string]interface{}
}

type WebsocketConnResp struct {
}

type MetricsResp struct {
	List []MetricsItem
}

type MetricsItem struct {
	Ud jwtx.UserData
}
