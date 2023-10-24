//go:build reqrep
// +build reqrep

package client

import (
	"fmt"
	"time"

	"github.com/learn-go/xmq/common"
	"github.com/learn-go/xmq/pkg"

	zmq "github.com/pebbe/zmq4"
)

type ReqRepClient struct{}

func init() {
	pkg.PkgClient = &ReqRepClient{}
}

func (s *ReqRepClient) Start() {
	//REQ 表示client端
	socket, _ := zmq.NewSocket(zmq.REQ)
	//Connect 请求建立连接，并指定传输层协议
	socket.Connect(fmt.Sprintf("tcp://127.0.0.1:%v", common.Port))
	fmt.Printf("client connect to server %v\n", common.Port)
	defer socket.Close()

	for i := 0; i < 10; i++ {
		//Send和Recv必须交替进行
		socket.Send("world", zmq.DONTWAIT) //非阻塞模式，异步发送（只是将数据写入本地buffer，并没有真正发送到网络上）
		resp, _ := socket.Recv(0)
		fmt.Printf("receive [%s]\n", resp)
		time.Sleep(5 * time.Second)
	}
}
