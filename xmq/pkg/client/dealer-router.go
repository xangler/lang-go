//go:build dealerrouter
// +build dealerrouter

package client

import (
	"fmt"
	"time"

	"github.com/learn-go/xmq/common"
	"github.com/learn-go/xmq/pkg"

	zmq "github.com/pebbe/zmq4"
)

type DealerRouterClient struct{}

func init() {
	pkg.PkgClient = &DealerRouterClient{}
}

func (s *DealerRouterClient) Start() {
	//DEALER 表示client端
	socket, _ := zmq.NewSocket(zmq.DEALER)
	//Connect 请求建立连接，并指定传输层协议
	socket.Connect(fmt.Sprintf("tcp://127.0.0.1:%v", common.Port))
	fmt.Printf("client connect to server %v\n", common.Port)
	defer socket.Close()

	for i := 0; i < 10; i++ {
		//Send和Recv没必要交替进行
		socket.Send("world", 0) //非阻塞模式，异步发送（只是将数据写入本地buffer，并没有真正发送到网络上）
		resp1, _ := socket.Recv(0)
		resp2, _ := socket.Recv(0)
		fmt.Printf("receive [%s %s]\n", resp1, resp2)
		time.Sleep(5 * time.Second)
	}
}
