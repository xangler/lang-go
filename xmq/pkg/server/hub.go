//go:build hub
// +build hub

package server

import (
	"encoding/base64"
	"fmt"

	"github.com/learn-go/xmq/common"
	"github.com/learn-go/xmq/pkg"

	zmq "github.com/pebbe/zmq4"
)

type HubServer struct{}

func init() {
	pkg.PkgServer = &HubServer{}
}

func (s *HubServer) Start() {
	//接收所有client的消息
	socket, _ := zmq.NewSocket(zmq.ROUTER)
	socket.Bind(fmt.Sprintf("tcp://127.0.0.1:%v", common.HubIn))
	fmt.Printf("server bind to port %d\n", common.HubIn)
	defer socket.Close()

	//把消息广播给所有client
	ctx, _ := zmq.NewContext()
	defer ctx.Term()
	publisher, _ := ctx.NewSocket(zmq.PUB)
	defer publisher.Close()
	publisher.Bind(fmt.Sprintf("tcp://127.0.0.1:%v", common.HubOut))
	fmt.Printf("pub bind to port %d\n", common.HubOut)

	for {
		//把接收到的client的消息再广播给所有client
		if addr, err := socket.RecvBytes(0); err == nil { //第一帧读出对端的地址
			client := base64.StdEncoding.EncodeToString(addr) //用对端地址来标识消息是谁发出来的
			if resp, err := socket.Recv(0); err == nil {
				if _, err := publisher.Send(client+"say: "+resp, 0); err != nil { //在消息前加上发送者的标识
					fmt.Println(err)
					break
				}
			} else {
				fmt.Println(err)
				break
			}
		} else {
			fmt.Println(err)
			break
		}
	}
}
