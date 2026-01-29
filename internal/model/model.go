package model

import "bifrost/common/jwtx"

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
	Ckv string `form:"ckv"`
}
