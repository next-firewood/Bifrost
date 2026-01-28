package config

import (
	"bifrost/common/jwtx"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Server struct {
		Port string `json:"port" yaml:"port"` // HTTP 端口
	} `json:"server" yaml:"server"`
	Auth jwtx.Auth `json:"auth" yaml:"auth"`

	vo.NacosClientParam `yaml:"nacosClientParam"`
	DataId              string `yaml:"dataId"`
	Group               string `yaml:"group"`
}

type Loader interface {
	Load() (*Config, error)
}

type LocalLoader struct {
	File string
}

func (l *LocalLoader) Load() (*Config, error) {
	data, err := os.ReadFile(l.File)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Load() (conf *Config, err error) {
	ct, err := clients.NewConfigClient(c.NacosClientParam)
	if err != nil {
		return conf, err
	}

	vcp := vo.ConfigParam{
		DataId: c.DataId,
		Group:  c.Group,
		OnChange: func(namespace, group, dataId, data string) {
			log.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	}

	content, err := ct.GetConfig(vcp)
	if err != nil {
		return conf, err
	}

	if content == "" {
		return conf, fmt.Errorf("nacos读取配置为空")
	}

	if err = ct.ListenConfig(vcp); err != nil {
		return conf, err
	}

	conf = &Config{}

	if err = yaml.Unmarshal([]byte(content), conf); err != nil {
		return conf, err
	}

	return conf, err
}
