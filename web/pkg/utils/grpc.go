package utils

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	GrpcConnBackoffMaxDelay = 10 // GrpcConnBackoffMaxDelay grpc conn reconnect max delay
	GrpcClientTimeout       = 15 // GrpcClientTimeout is timeout for calling grpc interface
)

func GetGrpcConnection(endpoint string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials())) // nolint
	// opts = append(opts, grpc.WithBackoffMaxDelay(GrpcConnBackoffMaxDelay*time.Second)) // nolint
	bc := backoff.DefaultConfig
	bc.MaxDelay = GrpcConnBackoffMaxDelay * time.Second
	p := grpc.ConnectParams{Backoff: bc}
	opts = append(opts, grpc.WithConnectParams(p), WithUnaryClientInterceptors(TimeoutInterceptor(GrpcClientTimeout*time.Second)))
	return grpc.Dial(endpoint, opts...)
}

// WithUnaryClientInterceptors uses given client unary interceptors.
func WithUnaryClientInterceptors(interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(interceptors...)
}

// TimeoutInterceptor is an interceptor that controls timeout.
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if timeout <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func ParseRpcInfo(info *grpc.UnaryServerInfo) (service, method string) {
	serMethods := strings.Split(info.FullMethod, "/")
	if len(serMethods) > 2 {
		method = serMethods[len(serMethods)-1]
		services := strings.Split(serMethods[1], ".")
		if len(services) > 1 {
			service = services[len(services)-1]
			return service, method
		}
	}
	return "", ""
}
