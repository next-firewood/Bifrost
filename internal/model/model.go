package model

import (
	"bifrost/common/jwtx"
)

type WebsocketConnReq struct {
	UserData map[string]interface{} `json:"userData"`
}

type WebsocketConnResp struct {
}

type MetricsResp struct {
	Total int64         `json:"total"`
	List  []MetricsItem `json:"list"`
}

type MetricsItem struct {
	Ud jwtx.UserData `json:"ud"`
}

type UserStatusReq struct {
	K string `form:"k"`
	V string `form:"v"`
}

type FilterHeader struct {
	FilterKey   string `header:"Filter-Key"`
	FilterValue string `header:"Filter-Value"`
}

type PushBroadcastReq struct {
	Msg string `json:"msg"`
}
