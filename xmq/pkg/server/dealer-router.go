//go:build dealerrouter
// +build dealerrouter

package server

import (
	"fmt"

	"github.com/learn-go/xmq/common"
	"github.com/learn-go/xmq/pkg"

	zmq "github.com/pebbe/zmq4"
)

type DealerRouterServer struct{}

func init() {
	pkg.PkgServer = &DealerRouterServer{}
}

func (s *DealerRouterServer) Start() {
	//ROUTER 表示server端
	socket, _ := zmq.NewSocket(zmq.ROUTER)
	//Bind 绑定端口，并指定传输层协议
	socket.Bind(fmt.Sprintf("tcp://127.0.0.1:%v", common.Port))
	fmt.Printf("server bind to port %d\n", common.Port)
	defer socket.Close()

	for {
		//Send和Recv没必要交替进行
		addr, _ := socket.RecvBytes(0) //接收到的第一帧表示对方的地址UUID
		resp, _ := socket.Recv(0)
		socket.SendBytes(addr, zmq.SNDMORE) //第一帧需要指明对方的地址，SNDMORE表示消息还没发完
		socket.Send("Hello", zmq.SNDMORE)   //如果不用SNDMORE表示这已经是最后一帧了，下一次Send就是下一段消息的第一帧了，需要指明对方的地址
		socket.Send(resp, 0)
	}
}
