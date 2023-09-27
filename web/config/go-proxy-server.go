package config

import (
	"github.com/learn-go/web/pkg/utils"
)

type ProxyConfig struct {
	HTTPPort   int    `yaml:"http_port"`
	RPCPort    int    `yaml:"rpc_port"`
	MetricPort int    `yaml:"metric_port"`
	ForwardURL string `yaml:"forward_url"`
}

func LoadProxyConfig(filePath string) (*ProxyConfig, error) {
	c := ProxyConfig{}
	if err := utils.Load(filePath, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
