package svc

import (
	"bifrost/common/wss"
	"bifrost/internal/config"
)

type ServerContext struct {
	Config *config.Config
	Hub    *wss.Hub
}

func NewServiceContext(c *config.Config) *ServerContext {
	return &ServerContext{
		Config: c,
		Hub:    wss.NewHub(),
	}
}
