package config

import (
	"fmt"
	"os"

	"github.com/learn-go/web/pkg/dbutils"
	"github.com/learn-go/web/pkg/utils"
)

type WebConfig struct {
	HTTPPort    int              `yaml:"http_port"`
	RPCPort     int              `yaml:"rpc_port"`
	MetricPort  int              `yaml:"metric_port"`
	CheckPeriod int              `yaml:"check_period"`
	Locker      *Locker          `yaml:"locker"`
	SQL         *dbutils.SQLPara `yaml:"sql"`
}

type Locker struct {
	Name     string `yaml:"name"`
	Period   int    `yaml:"period"`
	Dead     int64  `yaml:"dead"`
	WorkMark string `yaml:"work_mark"`
}

func LoadWebConfig(filePath string) (*WebConfig, error) {
	c := WebConfig{}
	if err := utils.Load(filePath, &c); err != nil {
		return nil, err
	}
	if v, ok := os.LookupEnv("MYSQL_ADDRESS"); ok {
		c.SQL.Addr = v
	}
	if v, ok := os.LookupEnv("MYSQL_DATABASE"); ok {
		c.SQL.DBName = v
	}
	if v, ok := os.LookupEnv("MYSQL_USERNAME"); ok {
		c.SQL.UserName = v
	}
	if v, ok := os.LookupEnv("MYSQL_PASSWORD"); ok {
		c.SQL.Password = v
	}
	c.Locker.WorkMark = fmt.Sprintf("%s_%s", c.Locker.WorkMark, utils.RandStringRunes(8))
	return &c, nil
}
