package leader

type ClientConfig struct {
	KubeConfig string `yaml:"kube_config"`
	MasterUrl  string `yaml:"master_url"`
	LockName   string `yaml:"lock_name"`
	Namespace  string `yaml:"namespace"`
}

type ForWardConfig struct {
	RPCPort int `yaml:"rpc_port"`
}

type PLFConfig struct {
	ForWardConfig  *ForWardConfig `yaml:"forward_config"`
	ClientConfig   *ClientConfig  `yaml:"client_config"`
	WithNotifyChan bool           `yaml:"with_notify_chan"`
	LeaderElect    bool           `yaml:"leader_elect"`
	ServerForward  bool           `yaml:"server_forward"`
}
