package common

const (
	Port   int32  = 30200
	HubIn  int32  = 30201
	HubOut int32  = 30202
	Prefix string = "demo"
)

type Server interface {
	Start()
}

type Client interface {
	Start()
}
