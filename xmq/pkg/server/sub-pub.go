//go:build subpub
// +build subpub

package server

import (
	"fmt"
	"time"

	"github.com/learn-go/xmq/common"
	"github.com/learn-go/xmq/pkg"

	zmq "github.com/pebbe/zmq4"
)

type SubPubServer struct{}

func init() {
	pkg.PkgServer = &SubPubServer{}
}

func (s *SubPubServer) Start() {
	ctx, _ := zmq.NewContext()
	defer ctx.Term()

	//PUB 表示publisher角色
	publisher, _ := ctx.NewSocket(zmq.PUB)
	defer publisher.Close()
	//Bind 绑定端口，并指定传输层协议
	publisher.Bind(fmt.Sprintf("tcp://127.0.0.1:%v", common.Port))
	fmt.Printf("server bind to port %d\n", common.Port)

	//publisher会把消息发送给所有subscriber，subscriber可以动态加入
	for i := 0; i < 5; i++ {
		//publisher只能调用send方法
		publisher.Send(common.Prefix+"Hello my followers", 0)
		publisher.Send(common.Prefix+"How are you", 0)
		fmt.Printf("loop %d send over\n", i+1)
		time.Sleep(10 * time.Second)
	}
	publisher.Send(common.Prefix+"END", 0)
}
