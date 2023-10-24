//go:build pushpull
// +build pushpull

package server

import (
	"fmt"
	"time"

	"github.com/learn-go/xmq/common"
	"github.com/learn-go/xmq/pkg"

	zmq "github.com/pebbe/zmq4"
)

type PushPullServer struct{}

func init() {
	pkg.PkgServer = &PushPullServer{}
}

func (s *PushPullServer) Start() {
	ctx, _ := zmq.NewContext()
	defer ctx.Term()

	//PUSH 表示pusher角色
	pusher, _ := ctx.NewSocket(zmq.PUSH)
	defer pusher.Close()
	//Bind 绑定端口，并指定传输层协议
	pusher.SetSndhwm(110)
	pusher.Bind(fmt.Sprintf("tcp://127.0.0.1:%v", common.Port))
	fmt.Printf("server bind to port %d\n", common.Port)

	//pusher把消息送给一个puller（采用公平轮转的方式选择一个puller）,puller可以动态加入
	for i := 0; i < 5; i++ {
		pusher.Send("Hello my followers", 0)
		pusher.Send("How are you", 0)
		fmt.Printf("loop %d send over\n", i+1)
		time.Sleep(5 * time.Second)
	}
	pusher.Send("END", 0)
}
