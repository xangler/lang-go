package leader

import (
	"errors"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var ErrBadAddress = errors.New("bad address")

type Forwarder struct {
	locker   sync.RWMutex
	opts     []grpc.DialOption
	lastAddr string
	lastConn *grpc.ClientConn
	conf     *ForWardConfig
}

func NewForwarder(conf *ForWardConfig, opts ...grpc.DialOption) *Forwarder {
	return &Forwarder{
		conf: conf,
		opts: opts,
	}
}

func (f *Forwarder) Close() error {
	f.locker.Lock()
	defer f.locker.Unlock()

	f.lastAddr = ""
	if f.lastConn == nil {
		return nil
	}
	return f.lastConn.Close()
}

func (f *Forwarder) ConnForAddr(addr string) (*grpc.ClientConn, error) {
	f.locker.Lock()
	defer f.locker.Unlock()

	if addr == "" {
		return nil, ErrBadAddress
	}
	if f.lastAddr == addr {
		return f.lastConn, nil
	}

	// grpc Dial never blocks, it is safe to have it in mutex
	fulAddr := fmt.Sprintf("%s:%d", addr, f.conf.RPCPort)
	log.Info("dial to ", fulAddr, " to forward request")
	conn, err := grpc.Dial(fulAddr, f.opts...)
	if err != nil {
		log.Error("dial leader candiate failed, ", err)
		return nil, err
	}

	if f.lastConn != nil {
		if err := f.lastConn.Close(); err != nil {
			log.Error("forwarder close old connection error: ", err)
		}
	}
	f.lastAddr = addr
	f.lastConn = conn
	return f.lastConn, nil
}

func (f *Forwarder) GetConnAddr() (*grpc.ClientConn, string) {
	f.locker.RLock()
	defer f.locker.RUnlock()
	return f.lastConn, f.lastAddr
}
