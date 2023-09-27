package config

import (
	"github.com/learn-go/web/pkg/leader"
	"github.com/learn-go/web/pkg/utils"
)

type LeaderServer struct {
	HTTPPort   int               `yaml:"http_port"`
	RPCPort    int               `yaml:"rpc_port"`
	MetricPort int               `yaml:"metric_port"`
	PLFConfig  *leader.PLFConfig `yaml:"plf_config"`
}

func LoadLeaderServer(filePath string) (*LeaderServer, error) {
	c := LeaderServer{}
	if err := utils.Load(filePath, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
