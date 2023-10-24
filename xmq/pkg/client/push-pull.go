//go:build pushpull
// +build pushpull

package client

import (
	"fmt"

	"github.com/learn-go/xmq/common"
	"github.com/learn-go/xmq/pkg"

	zmq "github.com/pebbe/zmq4"
)

type PushPullClient struct{}

func init() {
	pkg.PkgClient = &PushPullClient{}
}

func (s *PushPullClient) Start() {
	//PULL 表示puller角色
	puller, _ := zmq.NewSocket(zmq.PULL)
	defer puller.Close()

	//Bind 绑定端口，并指定传输层协议
	puller.Connect(fmt.Sprintf("tcp://127.0.0.1:%v", common.Port))
	fmt.Printf("client connect to server %v\n", common.Port)

	for {
		//接收广播
		if resp, err := puller.Recv(0); err == nil {
			fmt.Printf("receive [%s]\n", resp)
			if resp == "END" {
				break
			}
		} else {
			fmt.Println(err)
			break
		}
	}
}
