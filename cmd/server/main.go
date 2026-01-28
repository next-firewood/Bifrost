package main

import (
	"bifrost/internal/config"
	"bifrost/internal/route"
	"bifrost/svc"
	"github.com/gin-gonic/gin"
	"log"
)

var ConfigPath = "config.yaml"

func main() {
	loader := &config.LocalLoader{File: ConfigPath}

	cfg, err := loader.Load()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.DataId != "" && cfg.Group != "" {
		cfg, err = cfg.Load()
		if err != nil {
			log.Fatal(err)
		}
	}

	svcX := svc.NewServiceContext(cfg)

	go svcX.Hub.Run()

	r := gin.New()

	route.RegisterHTTP(r, svcX)

	r.Run(":" + cfg.Server.Port)
}
