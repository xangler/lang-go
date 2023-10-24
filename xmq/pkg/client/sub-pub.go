//go:build subpub
// +build subpub

package client

import (
	"fmt"

	"github.com/learn-go/xmq/common"
	"github.com/learn-go/xmq/pkg"

	zmq "github.com/pebbe/zmq4"
)

type SubPubClient struct{}

func init() {
	pkg.PkgClient = &SubPubClient{}
}

func (s *SubPubClient) Start() {
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	defer subscriber.Close()

	//Bind 绑定端口，并指定传输层协议
	subscriber.Connect(fmt.Sprintf("tcp://127.0.0.1:%v", common.Port))
	fmt.Printf("client connect to server %v\n", common.Port)
	subscriber.SetSubscribe(common.Prefix) //只接收前缀为prefix的消息

	for {
		//接收广播
		if resp, err := subscriber.Recv(0); err == nil {
			resp = resp[len(common.Prefix):] //去掉前缀
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
