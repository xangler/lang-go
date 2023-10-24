//go:build hub
// +build hub

package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/learn-go/xmq/common"
	"github.com/learn-go/xmq/pkg"

	zmq "github.com/pebbe/zmq4"
)

type HubClient struct{}

func init() {
	pkg.PkgClient = &HubClient{}
}

func (s *HubClient) Start() {
	//把消息广播给hub
	socket, _ := zmq.NewSocket(zmq.DEALER)
	socket.Connect(fmt.Sprintf("tcp://127.0.0.1:%v", common.HubIn))
	fmt.Printf("client connect to server %v\n", common.HubIn)
	defer socket.Close()

	//订阅hub的消息
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	defer subscriber.Close()
	subscriber.Connect(fmt.Sprintf("tcp://127.0.0.1:%v", common.HubOut))
	subscriber.SetSubscribe("")
	fmt.Printf("sub bind to port %d\n", common.HubOut)

	go func() {
		for {
			//把接收到的client的消息再广播给所有client
			if resp, err := subscriber.Recv(0); err == nil {
				fmt.Println(resp)
			} else {
				fmt.Println(err)
				break
			}
		}
	}()

	fmt.Println("please type message")
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		socket.Send(text, 0)
	}
}
