//go:build reqrep
// +build reqrep

package server

import (
	"fmt"

	"github.com/learn-go/xmq/common"
	"github.com/learn-go/xmq/pkg"

	zmq "github.com/pebbe/zmq4"
)

type ReqRepServer struct{}

func init() {
	pkg.PkgServer = &ReqRepServer{}
}

func (s *ReqRepServer) Start() {
	socket, _ := zmq.NewSocket(zmq.REP)
	//Bind 绑定端口，并指定传输层协议
	socket.Bind(fmt.Sprintf("tcp://127.0.0.1:%v", common.Port))
	fmt.Printf("server bind to port %d\n", common.Port)
	defer socket.Close()

	for {
		//Recv和Send必须交替进行
		resp, _ := socket.Recv(0)     //0表示阻塞模式
		socket.Send("Hello "+resp, 0) //同步发送
	}
}
