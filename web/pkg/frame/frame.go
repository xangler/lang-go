package frame

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"reflect"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Options struct{}

type RPCRegister struct {
	Register interface{}
	Server   interface{}
}

type HTTPRegister struct {
	Register interface{}
}

type ProtocServer struct {
	options       *Options
	rpcPort       int
	rpcRegisters  []*RPCRegister
	rpcSrv        *grpc.Server
	httpPort      int
	httpRegisters []*HTTPRegister
	httpSrv       *http.Server
	rpcOptions    []grpc.UnaryServerInterceptor
}

func NewProtocServer(opts *Options) *ProtocServer {
	return &ProtocServer{
		options:       opts,
		rpcRegisters:  []*RPCRegister{},
		httpRegisters: []*HTTPRegister{},
		rpcOptions:    []grpc.UnaryServerInterceptor{},
	}
}

func (s *ProtocServer) SetRpcOptions(rpcOptions []grpc.UnaryServerInterceptor) *ProtocServer {
	s.rpcOptions = rpcOptions
	return s
}

func (s *ProtocServer) SetRpcPort(port int) *ProtocServer {
	s.rpcPort = port
	return s
}

func (s *ProtocServer) RPC(register interface{}, server interface{}) *ProtocServer {
	s.rpcRegisters = append(s.rpcRegisters, &RPCRegister{
		Register: register,
		Server:   server,
	})
	return s
}

func (s *ProtocServer) SetHttpPort(port int) *ProtocServer {
	s.httpPort = port
	return s
}

func (s *ProtocServer) HTTP(register interface{}) *ProtocServer {
	s.httpRegisters = append(s.httpRegisters, &HTTPRegister{
		Register: register,
	})
	return s
}

func (s *ProtocServer) Run() {
	go s.processRPC()
	go s.processHTTP()
}

func (s *ProtocServer) Close() {
	if s.httpSrv != nil {
		s.httpSrv.Close()
		s.httpSrv = nil
	}

	if s.rpcSrv != nil {
		s.rpcSrv.Stop()
		s.rpcSrv = nil
	}
}

func (s *ProtocServer) processRPC() {
	if len(s.rpcRegisters) == 0 {
		log.Panic("grpc failed: no register function")
	}

	if err := s.checkRPC(); err != nil {
		log.Panic(err.Error())
	}

	if err := s.runRPC(); err != nil {
		log.Panic(err.Error())
	}
}

func (s *ProtocServer) checkRPC() error {
	for i, rpcRegister := range s.rpcRegisters {
		mt := reflect.TypeOf(rpcRegister.Register)
		if mt.Kind() != reflect.Func {
			return fmt.Errorf("the [%v]th rpc register must be function", i)
		}

		// check register must be like func RegisterXXXXServer(s *grpc.Server, srv XXXXServer)
		in := mt.NumIn()
		out := mt.NumOut()

		if in != 2 || out != 0 {
			return fmt.Errorf("the number of inputs or outputs of the [%v]th rpc register is not correct", i)
		}

		msgFmt := "the definition of the [%v]th rpc register is not correct: %v"

		for i := 0; i < in; i++ {
			arg := mt.In(i)

			switch i {
			case 0:
				if _, ok := arg.(grpc.ServiceRegistrar); ok {
					return fmt.Errorf(msgFmt, i, "1st arg is not implemtent \"grpc.ServiceRegistrar\" func.")
				}
			case 1:
				if !reflect.TypeOf(rpcRegister.Server).Implements(arg) {
					return fmt.Errorf(msgFmt, i, "2nd arg is not \"XXXXServer\"")
				}
			}
		}
	}

	return nil
}

func (s *ProtocServer) runRPC() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", s.rpcPort))
	if err != nil {
		return err
	}
	interceptors := grpc.ChainUnaryInterceptor(s.rpcOptions...)
	opts := []grpc.ServerOption{
		interceptors,
	}
	rpcSrv := grpc.NewServer(opts...)
	s.rpcSrv = rpcSrv

	for _, rpcRegister := range s.rpcRegisters {
		register := reflect.ValueOf(rpcRegister.Register)
		args := make([]reflect.Value, 2)
		args[0] = reflect.ValueOf(rpcSrv)
		args[1] = reflect.ValueOf(rpcRegister.Server)
		register.Call(args)
	}

	log.Infof("starting rpc server on port [%v]\n", s.rpcPort)

	return rpcSrv.Serve(l)
}

func (s *ProtocServer) processHTTP() {
	if len(s.httpRegisters) != 0 {
		if err := s.checkHTTP(); err != nil {
			log.Panic(err.Error())
		}

		if err := s.runHTTP(); err != nil {
			if err != http.ErrServerClosed {
				log.Panic(err.Error())
			}
		}
	}
}

func (s *ProtocServer) checkHTTP() error {
	for i, httpRegister := range s.httpRegisters {
		mt := reflect.TypeOf(httpRegister.Register)
		if mt.Kind() != reflect.Func {
			return fmt.Errorf("the [%v]th http register must be function", i)
		}

		// check register must be like func RegisterXXXXHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
		in := mt.NumIn()
		out := mt.NumOut()

		if in != 4 || out != 1 {
			return fmt.Errorf("the number of inputs or outputs of the [%v]th http register is not correct", i)
		}

		msgFmt := "the definition of the [%v]th http register is not correct: %v"

		for i := 0; i < in; i++ {
			arg := mt.In(i)

			switch i {
			case 0:
				if !arg.Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
					return fmt.Errorf(msgFmt, i, "1st arg is not \"context.Context\"")
				}
			case 1:
				if arg != reflect.TypeOf(&runtime.ServeMux{}) {
					return fmt.Errorf(msgFmt, i, "2nd arg is not \"*runtime.ServeMux\"")
				}
			case 2:
				if arg.Kind() != reflect.String {
					return fmt.Errorf(msgFmt, i, "3rd arg is not \"string\"")
				}
			case 3:
				if arg.Kind() != reflect.Slice &&
					!arg.Elem().Implements(reflect.TypeOf((*grpc.DialOption)(nil)).Elem()) {
					return fmt.Errorf(msgFmt, i, "4th arg is not \"[]grpc.DialOption\"")
				}
			}
		}

		if !mt.Out(0).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			return fmt.Errorf(msgFmt, i, "output is not \"error\"")
		}
	}

	return nil
}

func (s *ProtocServer) runHTTP() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	for _, httpRegister := range s.httpRegisters {
		reg := reflect.ValueOf(httpRegister.Register)
		args := make([]reflect.Value, 4)
		args[0] = reflect.ValueOf(ctx)
		args[1] = reflect.ValueOf(gwmux)
		args[2] = reflect.ValueOf(fmt.Sprintf("localhost:%v", s.rpcPort))
		args[3] = reflect.ValueOf(opts)
		res := reg.Call(args)

		err := res[0].Interface()
		if err != nil {
			return err.(error)
		}
	}

	// Start HTTP rpcsrv (and proxy calls to gRPC rpcsrv endpoint)
	log.Infof("starting http server on port [%v]", s.httpPort)

	httpSrv := &http.Server{Addr: fmt.Sprintf(":%v", s.httpPort), Handler: gwmux}
	s.httpSrv = httpSrv

	return httpSrv.ListenAndServe()
}
